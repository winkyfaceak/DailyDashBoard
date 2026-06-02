package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomRecipeUsesAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected JSON accept header, got %q", got)
		}

		fmt.Fprint(w, `{"meals":[{"idMeal":"52772","strMeal":"Teriyaki Chicken Casserole","strCategory":"Chicken","strArea":"Japanese","strInstructions":"Bake until bubbling.","strMealThumb":"https://example.com/meal.jpg","strIngredient1":"Chicken","strMeasure1":"3 breasts"}]}`)
	}))
	defer server.Close()

	restoreRandomRecipeURL(t, server.URL)

	recipe, err := RandomRecipe()
	if err != nil {
		t.Fatal(err)
	}

	if recipe.Name != "Teriyaki Chicken Casserole" {
		t.Fatalf("expected recipe name, got %q", recipe.Name)
	}

	if recipe.ImageURL != "https://example.com/meal.jpg" {
		t.Fatalf("expected meal image URL, got %q", recipe.ImageURL)
	}
}

func TestDecodeRandomRecipeCleansIngredientsAndText(t *testing.T) {
	recipe, err := decodeRandomRecipe(strings.NewReader(`{
		"meals": [{
			"idMeal": "52883",
			"strMeal": " Sticky&nbsp;Toffee Pudding ",
			"strCategory": " Dessert ",
			"strArea": " British ",
			"strInstructions": "Bake\r\nuntil risen.\n\nServe warm.",
			"strMealThumb": " https://example.com/sticky.jpg ",
			"strTags": "Pudding, Cake, ",
			"strSource": " https://example.com/source ",
			"strYoutube": " https://example.com/video ",
			"strIngredient1": " Butter ",
			"strMeasure1": " 50g ",
			"strIngredient2": "",
			"strMeasure2": "",
			"strIngredient3": " Cream ",
			"strMeasure3": ""
		}]
	}`))
	if err != nil {
		t.Fatal(err)
	}

	if recipe.Name != "Sticky Toffee Pudding" {
		t.Fatalf("expected cleaned name, got %q", recipe.Name)
	}

	if recipe.Instructions != "Bake until risen. Serve warm." {
		t.Fatalf("expected cleaned instructions, got %q", recipe.Instructions)
	}

	if len(recipe.Ingredients) != 2 {
		t.Fatalf("expected 2 ingredients, got %d", len(recipe.Ingredients))
	}

	if recipe.Ingredients[0].Name != "Butter" || recipe.Ingredients[0].Measure != "50g" {
		t.Fatalf("expected cleaned first ingredient, got %#v", recipe.Ingredients[0])
	}

	if recipe.Ingredients[1].Name != "Cream" || recipe.Ingredients[1].Measure != "" {
		t.Fatalf("expected ingredient without measure, got %#v", recipe.Ingredients[1])
	}

	if strings.Join(recipe.Tags, ",") != "Pudding,Cake" {
		t.Fatalf("expected cleaned tags, got %#v", recipe.Tags)
	}
}

func TestDecodeRandomRecipeRejectsEmptyMeals(t *testing.T) {
	_, err := decodeRandomRecipe(strings.NewReader(`{"meals":[]}`))
	if err == nil {
		t.Fatal("expected empty meals error")
	}
}

func restoreRandomRecipeURL(t *testing.T, testURL string) {
	t.Helper()

	originalURL := randomRecipeURL
	randomRecipeURL = testURL

	t.Cleanup(func() {
		randomRecipeURL = originalURL
	})
}
