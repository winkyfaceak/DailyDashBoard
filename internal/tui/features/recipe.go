package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	recipeTextWidth = 66
	recipeImageCols = recipeTextWidth
	recipeImageRows = 18
)

// RandomRecipe shows a random recipe from TheMealDB.
func RandomRecipe() *Feature {
	return &Feature{
		MenuLabel:       "random recipe",
		Loading:         "Finding a random recipe...",
		ErrLabel:        "Recipe",
		BoxWidth:        84,
		Accent:          lipgloss.Color("#D98E4A"),
		BackgroundWords: []string{"simmer", "whisk", "spice", "recipe", "❦"},
		Fetch: func() (Session, error) {
			recipe, err := daily.RandomRecipe()
			if err != nil {
				return nil, err
			}

			return &RecipeSession{
				Recipe:     recipe,
				imageState: newImageState(renderRecipeImage(recipe)),
			}, nil
		},
	}
}

// RecipeSession renders a recipe with its food image and source links.
type RecipeSession struct {
	imageState
	Recipe daily.Recipe
}

func (s *RecipeSession) View() string {
	return FormatRecipe(s.Recipe, s.imageText("Rendering food image..."))
}

func (s *RecipeSession) Links() []LinkAction {
	var links []LinkAction

	if s.Recipe.ImageURL != "" {
		links = append(links, LinkAction{
			Key:     "o",
			URL:     s.Recipe.ImageURL,
			Help:    "o to open food image",
			Message: "Opening food image...",
		})
	}

	if s.Recipe.SourceURL != "" {
		links = append(links, LinkAction{
			Key:     "r",
			URL:     s.Recipe.SourceURL,
			Help:    "r to open recipe",
			Message: "Opening recipe source...",
		})
	}

	if s.Recipe.YouTubeURL != "" {
		links = append(links, LinkAction{
			Key:     "v",
			URL:     s.Recipe.YouTubeURL,
			Help:    "v to open video",
			Message: "Opening cooking video...",
		})
	}

	return links
}

func renderRecipeImage(recipe daily.Recipe) terminalimage.Image {
	if recipe.ImageURL == "" {
		return terminalimage.Image{Text: subtleTextStyle.Render("Food image unavailable.")}
	}

	renderedImage, err := terminalimage.RenderFromURL(recipe.ImageURL, recipeImageCols, recipeImageRows)
	if err != nil {
		return terminalimage.Image{Text: subtleTextStyle.Render(fmt.Sprintf("Food image unavailable: %v", err))}
	}

	return renderedImage
}

func FormatRecipe(recipe daily.Recipe, renderedImage string) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(recipeTextWidth).Render(formatRecipeHeader(recipe)))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n")
	}

	builder.WriteString(labelStyle.Render("Recipe: "))
	builder.WriteString(strongTextStyle.Render(recipe.Name))
	builder.WriteString("\n")

	meta := formatRecipeMeta(recipe)
	if meta != "" {
		builder.WriteString(subtleTextStyle.Render(meta))
		builder.WriteString("\n")
	}

	if len(recipe.Ingredients) > 0 {
		builder.WriteString("\n")
		builder.WriteString(labelStyle.Render("Ingredients:"))
		builder.WriteString("\n")
		builder.WriteString(formatRecipeIngredients(recipe.Ingredients))
		builder.WriteString("\n")
	}

	if recipe.Instructions != "" {
		builder.WriteString("\n")
		builder.WriteString(labelStyle.Render("Instructions:"))
		builder.WriteString("\n")
		builder.WriteString(textutil.Wrap(textutil.Truncate(recipe.Instructions, 850), recipeTextWidth))
		builder.WriteString("\n")
	}

	if recipe.SourceURL != "" {
		builder.WriteString("\n")
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(recipe.SourceURL).Render("press r to open recipe source"))
	}

	if recipe.YouTubeURL != "" {
		builder.WriteString("\n")
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(recipe.YouTubeURL).Render("press v to open cooking video"))
	}

	builder.WriteString("\n")
	builder.WriteString(subtleTextStyle.Render("Recipe and image provided courtesy of TheMealDB."))

	return builder.String()
}

func formatRecipeHeader(recipe daily.Recipe) string {
	meta := formatRecipeOrigin(recipe)
	if meta == "" {
		return "Random Recipe"
	}

	return "Random Recipe -> " + meta
}

func formatRecipeMeta(recipe daily.Recipe) string {
	var parts []string

	if recipe.Category != "" {
		parts = append(parts, recipe.Category)
	}

	if recipe.Area != "" {
		parts = append(parts, recipe.Area)
	}

	if len(recipe.Tags) > 0 {
		parts = append(parts, strings.Join(textutil.LimitStrings(recipe.Tags, 4), ", "))
	}

	return strings.Join(parts, " • ")
}

func formatRecipeOrigin(recipe daily.Recipe) string {
	var parts []string

	if recipe.Area != "" {
		parts = append(parts, recipe.Area)
	}

	if recipe.Category != "" {
		parts = append(parts, recipe.Category)
	}

	return strings.Join(parts, " ")
}

func formatRecipeIngredients(ingredients []daily.RecipeIngredient) string {
	limit := 12
	if len(ingredients) < limit {
		limit = len(ingredients)
	}

	var lines []string
	for _, ingredient := range ingredients[:limit] {
		line := recipeIngredientLine(ingredient)
		if line == "" {
			continue
		}

		lines = append(lines, textutil.IndentWrapped("- "+line, recipeTextWidth, "  "))
	}

	if remaining := len(ingredients) - limit; remaining > 0 {
		lines = append(lines, subtleTextStyle.Render(fmt.Sprintf("  ...and %d more ingredients", remaining)))
	}

	return strings.Join(lines, "\n")
}

func recipeIngredientLine(ingredient daily.RecipeIngredient) string {
	if ingredient.Measure == "" {
		return ingredient.Name
	}

	return ingredient.Measure + " " + ingredient.Name
}
