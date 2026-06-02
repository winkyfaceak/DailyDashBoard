package api

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

var randomRecipeURL = "https://www.themealdb.com/api/json/v1/1/random.php"

type Recipe struct {
	ID           string
	Name         string
	Category     string
	Area         string
	Instructions string
	ImageURL     string
	SourceURL    string
	YouTubeURL   string
	Tags         []string
	Ingredients  []RecipeIngredient
}

type RecipeIngredient struct {
	Name    string
	Measure string
}

type randomRecipeResponse struct {
	Meals []map[string]any `json:"meals"`
}

func RandomRecipe() (Recipe, error) {
	req, err := http.NewRequest(http.MethodGet, randomRecipeURL, nil)
	if err != nil {
		return Recipe{}, fmt.Errorf("error creating recipe request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return Recipe{}, fmt.Errorf("error sending recipe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Recipe{}, fmt.Errorf("recipe API error: %s", resp.Status)
	}

	recipe, err := decodeRandomRecipe(resp.Body)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func decodeRandomRecipe(reader io.Reader) (Recipe, error) {
	var response randomRecipeResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return Recipe{}, fmt.Errorf("error decoding recipe response: %w", err)
	}

	if len(response.Meals) == 0 {
		return Recipe{}, fmt.Errorf("recipe API returned no meals")
	}

	recipe := recipeFromMeal(response.Meals[0])
	if recipe.Name == "" {
		return Recipe{}, fmt.Errorf("recipe API returned no meal name")
	}

	if recipe.Instructions == "" && len(recipe.Ingredients) == 0 {
		return Recipe{}, fmt.Errorf("recipe API returned no recipe details")
	}

	return recipe, nil
}

func recipeFromMeal(meal map[string]any) Recipe {
	return Recipe{
		ID:           mealString(meal, "idMeal"),
		Name:         mealString(meal, "strMeal"),
		Category:     mealString(meal, "strCategory"),
		Area:         mealString(meal, "strArea"),
		Instructions: mealString(meal, "strInstructions"),
		ImageURL:     mealString(meal, "strMealThumb"),
		SourceURL:    mealString(meal, "strSource"),
		YouTubeURL:   mealString(meal, "strYoutube"),
		Tags:         recipeTags(mealString(meal, "strTags")),
		Ingredients:  recipeIngredients(meal),
	}
}

func recipeIngredients(meal map[string]any) []RecipeIngredient {
	var ingredients []RecipeIngredient

	for i := 1; i <= 20; i++ {
		ingredient := mealString(meal, fmt.Sprintf("strIngredient%d", i))
		if ingredient == "" {
			continue
		}

		ingredients = append(ingredients, RecipeIngredient{
			Name:    ingredient,
			Measure: mealString(meal, fmt.Sprintf("strMeasure%d", i)),
		})
	}

	return ingredients
}

func recipeTags(value string) []string {
	var tags []string

	for _, tag := range strings.Split(value, ",") {
		tag = cleanRecipeText(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

func mealString(meal map[string]any, key string) string {
	value, ok := meal[key]
	if !ok || value == nil {
		return ""
	}

	text, ok := value.(string)
	if !ok {
		text = fmt.Sprint(value)
	}

	return cleanRecipeText(text)
}

func cleanRecipeText(value string) string {
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}
