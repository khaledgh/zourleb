package models

// GetLocale implements i18n.Translatable for each translation entity so the
// generic resolver can pick the best locale row.

func (t TourTranslation) GetLocale() string     { return t.Locale }
func (t AgencyTranslation) GetLocale() string   { return t.Locale }
func (t CategoryTranslation) GetLocale() string { return t.Locale }
func (t RegionTranslation) GetLocale() string   { return t.Locale }
func (t ProductTranslation) GetLocale() string  { return t.Locale }
