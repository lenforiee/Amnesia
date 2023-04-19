package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/lenforiee/AmnesiaGUI/bundles"
)

type ClassicTheme struct {
	variant fyne.ThemeVariant

	// theming purposes.
	defaultTheme fyne.Theme
}

func ClassicThemeDark() *ClassicTheme {
	defaultTheme := theme.DefaultTheme()
	return &ClassicTheme{variant: theme.VariantDark, defaultTheme: defaultTheme}
}

func ClassicThemeLight() *ClassicTheme {
	defaultTheme := theme.DefaultTheme()
	return &ClassicTheme{variant: theme.VariantLight, defaultTheme: defaultTheme}
}

func (t *ClassicTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return t.defaultTheme.Font(s)
	}
	if s.Bold {
		if s.Italic {
			return t.defaultTheme.Font(s)
		}
		return bundles.FontJetBrainsMonoBoldTtf
	}
	if s.Italic {
		return t.defaultTheme.Font(s)
	}
	return bundles.FontJetBrainsMonoMediumTtf
}

func (t *ClassicTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return t.defaultTheme.Color(n, t.variant)
}

func (t *ClassicTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return t.defaultTheme.Icon(n)
}

func (t *ClassicTheme) Size(n fyne.ThemeSizeName) float32 {
	return t.defaultTheme.Size(n)
}
