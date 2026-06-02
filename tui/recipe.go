package tui

import (
	"DailyDashBoard/dailys/api"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const (
	recipeTextWidth = 66
	recipeImageCols = recipeTextWidth
	recipeImageRows = 18
)

type recipeResultMsg struct {
	Recipe api.Recipe
	Image  terminalImage
	Err    error
}

type recipeImageReadyMsg struct{}

func fetchRecipeCmd() tea.Cmd {
	return func() tea.Msg {
		recipe, err := api.RandomRecipe()
		if err != nil {
			return recipeResultMsg{Err: err}
		}

		return recipeResultMsg{
			Recipe: recipe,
			Image:  renderRecipeImage(recipe),
		}
	}
}

func renderRecipeImage(recipe api.Recipe) terminalImage {
	if recipe.ImageURL == "" {
		return terminalImage{Text: subtleTextStyle.Render("Food image unavailable.")}
	}

	renderedImage, err := renderTerminalImageFromURL(recipe.ImageURL, recipeImageCols, recipeImageRows)
	if err != nil {
		return terminalImage{Text: subtleTextStyle.Render(fmt.Sprintf("Food image unavailable: %v", err))}
	}

	return renderedImage
}

func formatRecipe(recipe api.Recipe, renderedImage string) string {
	var builder strings.Builder

	builder.WriteString(recipeHeaderStyle.Render(formatRecipeHeader(recipe)))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n")
	}

	builder.WriteString(labelStyle.Render("Recipe: "))
	builder.WriteString(recipeNameStyle.Render(recipe.Name))
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
		builder.WriteString(wrapText(truncateText(recipe.Instructions, 850), recipeTextWidth))
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

func formatRecipeHeader(recipe api.Recipe) string {
	meta := formatRecipeOrigin(recipe)
	if meta == "" {
		return "Random Recipe"
	}

	return "Random Recipe -> " + meta
}

func formatRecipeMeta(recipe api.Recipe) string {
	var parts []string

	if recipe.Category != "" {
		parts = append(parts, recipe.Category)
	}

	if recipe.Area != "" {
		parts = append(parts, recipe.Area)
	}

	if len(recipe.Tags) > 0 {
		parts = append(parts, strings.Join(limitStrings(recipe.Tags, 4), ", "))
	}

	return strings.Join(parts, " • ")
}

func formatRecipeOrigin(recipe api.Recipe) string {
	var parts []string

	if recipe.Area != "" {
		parts = append(parts, recipe.Area)
	}

	if recipe.Category != "" {
		parts = append(parts, recipe.Category)
	}

	return strings.Join(parts, " ")
}

func formatRecipeIngredients(ingredients []api.RecipeIngredient) string {
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

		lines = append(lines, indentWrappedText("- "+line, recipeTextWidth, "  "))
	}

	if remaining := len(ingredients) - limit; remaining > 0 {
		lines = append(lines, subtleTextStyle.Render(fmt.Sprintf("  ...and %d more ingredients", remaining)))
	}

	return strings.Join(lines, "\n")
}

func recipeIngredientLine(ingredient api.RecipeIngredient) string {
	if ingredient.Measure == "" {
		return ingredient.Name
	}

	return ingredient.Measure + " " + ingredient.Name
}
