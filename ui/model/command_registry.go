package model

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/bjarneo/cliamp/provider"
	"github.com/bjarneo/cliamp/ui"
)

// commandMode identifies the UI contexts in which a command is available.
// A command may belong to more than one context.
type commandMode uint64

const (
	commandModeMain commandMode = 1 << iota
	commandModeProvider
	commandModeEQ
	commandModeVolume
	commandModeShuffle
	commandModeRepeat
	commandModeSpeed
	commandModeProviderPill
	commandModeKeymap
	commandModeKeymapSearch
	commandModeFileBrowser
	commandModeFileBrowserSearch
	commandModeNavBrowser
	commandModeNavSearch
	commandModePlaylistManager
	commandModePlaylistManagerInput
	commandModePlaylistManagerDirs
	commandModePlaylistPicker
	commandModePlaylistPickerInput
	commandModeQueue
	commandModeSearch
	commandModeNetSearch
	commandModeSpotSearch
	commandModeJump
	commandModeURL
	commandModeLyrics
	commandModeThemePicker
	commandModeVisPicker
	commandModeDevicePicker
	commandModeInfo
	commandModeThemePickerFilter
	commandModeVisPickerFilter
	commandModeProviderSearch
	commandModeSubs
	commandModeSubsFilter
)

const commandModeAny = ^commandMode(0)

// commandSpec is the single source of metadata for in-app commands. Dispatch
// remains in the focused handlers, while keymap, plugin reservations, and help
// all consume this description.
type commandSpec struct {
	Mode        commandMode
	Keys        []string // Bubbletea KeyPressMsg.String values.
	KeyLabel    string   // Human-readable key label.
	Label       string
	LabelFor    func(Model) string
	Enabled     func(Model) bool
	Destructive bool
	Keymap      bool
	ContextHelp bool
	Prominent   bool
	Primary     bool
	Cancel      bool
	Help        bool
}

func (c commandSpec) enabled(m Model) bool {
	return c.Enabled == nil || c.Enabled(m)
}

func (c commandSpec) label(m Model) string {
	if c.LabelFor != nil {
		return c.LabelFor(m)
	}
	return c.Label
}

// commandRegistry deliberately lists every core-reserved key, including text
// editor keys that are not shown in the global keymap. Keep key labels in this
// table so the keymap cannot drift from plugin key reservations.
var commandRegistry = []commandSpec{
	{Mode: commandModeMain | commandModeEQ | commandModeSpeed, Keys: []string{"space"}, KeyLabel: "Space", Label: "Přehrát / pauza", Keymap: true, ContextHelp: true, Primary: true},
	{Mode: commandModeMain, Keys: []string{"s"}, KeyLabel: "s", Label: "Stop", Keymap: true},
	{Mode: commandModeMain, Keys: []string{">", "."}, KeyLabel: "> .", Label: "Další", Keymap: true, ContextHelp: true, Prominent: true},
	{Mode: commandModeMain, Keys: []string{"<", ","}, KeyLabel: "< ,", Label: "Předchozí", Keymap: true, ContextHelp: true, Prominent: true},
	{Mode: commandModeMain, Keys: []string{"left", "right"}, KeyLabel: "Left Right", Label: "Posun ±5 s", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"shift+left", "shift+right"}, KeyLabel: "Shift+Left Right", Label: "Velký posun", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "j"}, KeyLabel: "Nj", Label: "Skok na N×10 % skladby (např. 7j = 70 %)", Keymap: true},
	{Mode: commandModeMain | commandModeVolume | commandModeShuffle | commandModeRepeat, Keys: []string{"+", "=", "-"}, KeyLabel: "+ -", Label: "Hlasitost", Keymap: true},
	{Mode: commandModeMain | commandModeSpeed, Keys: []string{"]", "["}, KeyLabel: "] [", Label: "Rychlost ±0,25×", Keymap: true},
	{Mode: commandModeMain | commandModeVolume | commandModeShuffle | commandModeRepeat, Keys: []string{"z"}, KeyLabel: "z", Label: "Náhodně", Keymap: true},
	{Mode: commandModeMain | commandModeVolume | commandModeShuffle | commandModeRepeat, Keys: []string{"r"}, KeyLabel: "r", Label: "Opakování", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"m"}, KeyLabel: "m", Label: "Mono", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"e"}, KeyLabel: "e", Label: "Další předvolba EQ", Enabled: func(m Model) bool { return !m.simplified }, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"t"}, KeyLabel: "t", Label: "Vzhled", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"v"}, KeyLabel: "v", Label: "Další vizualizace", Enabled: func(m Model) bool { return !m.simplified }, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+v"}, KeyLabel: "Ctrl+V", Label: "Vybrat vizualizaci", Enabled: func(m Model) bool { return !m.simplified }, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"V"}, KeyLabel: "V", Label: "Vizualizace přes celou obrazovku", Enabled: func(m Model) bool { return !m.simplified }, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"up", "down", "k", "j"}, KeyLabel: "Up Down", Label: "Pohyb v seznamu / úprava EQ", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"pgup", "pgdown", "ctrl+u", "ctrl+d"}, KeyLabel: "PgUp PgDn / Ctrl+U D", Label: "Posun o stránku", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"home", "end", "g", "G"}, KeyLabel: "Home End / g G", Label: "Na začátek / konec seznamu", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"shift+up", "shift+down"}, KeyLabel: "Shift+Up Down", Label: "Posunout skladbu nahoru/dolů", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"h", "l"}, KeyLabel: "h l", Label: "Kurzor EQ vlevo/vpravo", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Přehrát vybranou", Keymap: true, ContextHelp: true, Primary: true},
	{Mode: commandModeMain, Keys: []string{"f"}, KeyLabel: "f", Label: "Záložka / oblíbená stanice", LabelFor: func(m Model) string {
		if m.selectedPlaylistStarAction() == starRadioFavorite {
			return "Oblíbená stanice"
		}
		return "Záložka"
	}, Enabled: func(m Model) bool { return m.selectedPlaylistStarAction() != starUnavailable }, Keymap: true, ContextHelp: true, Prominent: true},
	{Mode: commandModeMain, Keys: []string{"n"}, KeyLabel: "n", Label: "Oblíbená skladba", Enabled: func(m Model) bool {
		return m.focus == focusPlaylist && m.playlist != nil && m.favMgr != nil && m.plCursor >= 0 && m.plCursor < m.playlist.Len()
	}, Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"a"}, KeyLabel: "a", Label: "Zařadit jako další", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"A"}, KeyLabel: "A", Label: "Fronta", Keymap: true},
	{Mode: commandModeMain | commandModeProvider, Keys: []string{"F"}, KeyLabel: "F", Label: "Odebírané pořady", Enabled: func(m Model) bool { return m.hasSubscriptions() }, Keymap: true},

	{Mode: commandModeSubs, Keys: []string{"up", "down", "j", "k"}, KeyLabel: "Up Down", Label: "Pohyb", Keymap: true},
	{Mode: commandModeSubs, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr", Keymap: true, ContextHelp: true},
	{Mode: commandModeSubs, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Přidat díly a přehrát první", Keymap: true, ContextHelp: true, Primary: true},
	{Mode: commandModeSubs, Keys: []string{"a"}, KeyLabel: "a", Label: "Přidat díly", Keymap: true, ContextHelp: true},
	{Mode: commandModeSubs, Keys: []string{"q"}, KeyLabel: "q", Label: "Přidat a zařadit díly", Keymap: true},
	{Mode: commandModeSubs, Keys: []string{"l"}, KeyLabel: "l", Label: "Nejnovější díl do fronty", Keymap: true, ContextHelp: true},
	{Mode: commandModeSubs, Keys: []string{"L"}, KeyLabel: "L", Label: "Nejnovější z každého pořadu", Keymap: true, ContextHelp: true},
	{Mode: commandModeSubs, Keys: []string{"esc", "F"}, KeyLabel: "Esc", Label: "Zavřít", Keymap: true},
	{Mode: commandModeSubsFilter, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Použít filtr", Keymap: true, Primary: true},
	{Mode: commandModeSubsFilter, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zrušit filtr", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"x"}, KeyLabel: "x", Label: "Odebrat skladbu ze seznamu", Destructive: true, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"w"}, KeyLabel: "w", Label: "Uložit výběr do playlistu", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"o"}, KeyLabel: "o", Label: "Procházet soubory", Keymap: true},
	{Mode: commandModeMain | commandModeProvider, Keys: []string{"N"}, KeyLabel: "N", Label: "Playlisty", LabelFor: func(m Model) string {
		if _, ok := m.selectedTrackArtistBrowserTarget(); ok {
			return "Interpret"
		}
		return "Playlisty"
	}, Enabled: func(m Model) bool { return m.canOpenProviderBrowser() }, Keymap: true, ContextHelp: true, Prominent: true},
	{Mode: commandModeMain, Keys: []string{"L"}, KeyLabel: "L", Label: "Místní playlisty", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"R"}, KeyLabel: "R", Label: "Rádia", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"O"}, KeyLabel: "O", Label: "Podcasty", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"S"}, KeyLabel: "S", Label: "Spotify", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"P"}, KeyLabel: "P", Label: "Plex", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"Y"}, KeyLabel: "Y", Label: "YouTube", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"C"}, KeyLabel: "C", Label: "SoundCloud", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"X"}, KeyLabel: "X", Label: "Mixcloud", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"M"}, KeyLabel: "M", Label: "NetEase", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"J"}, KeyLabel: "J", Label: "Jellyfin", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"E"}, KeyLabel: "E", Label: "Emby", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"B"}, KeyLabel: "B", Label: "Audiobookshelf", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"Q"}, KeyLabel: "Q", Label: "Qobuz", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"T"}, KeyLabel: "T", Label: "Tidal", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+j"}, KeyLabel: "Ctrl+J", Label: "Skok na čas", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"p"}, KeyLabel: "p", Label: "Správa playlistů", Keymap: true},
	{Mode: commandModeProvider, Keys: []string{"p"}, KeyLabel: "p", Label: "Správa playlistů", Keymap: true, ContextHelp: true, Enabled: func(m Model) bool {
		return m.isActiveProvider("Local") && m.localProvider != nil
	}},
	{Mode: commandModeMain, Keys: []string{"ctrl+h"}, KeyLabel: "Ctrl+H", Label: "Záhlaví alb", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+g"}, KeyLabel: "Ctrl+G", Label: "Řádek nápovědy", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+b"}, KeyLabel: "Ctrl+B", Label: "Panel nastavení", Enabled: func(m Model) bool {
		return !m.simplified && m.layout.tier == layoutFull
	}, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"i"}, KeyLabel: "i", Label: "Info o skladbě", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain | commandModeInfo, Keys: []string{"ctrl+i"}, KeyLabel: "Ctrl+I", Label: "Metadata", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+s"}, KeyLabel: "Ctrl+S", Label: "Stáhnout skladbu", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+x"}, KeyLabel: "Ctrl+X", Label: "Rozbalit / sbalit", Enabled: func(m Model) bool { return !m.simplified }, Keymap: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+x"}, KeyLabel: "Ctrl+X", Label: "Rozbalit", Enabled: func(m Model) bool {
		return !m.simplified && !m.heightExpanded && m.layout.bodyRows > m.plVisible
	}},
	{Mode: commandModeMain, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr seznamu", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"ctrl+f"}, KeyLabel: "Ctrl+F", Label: "Hledat", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"u"}, KeyLabel: "u", Label: "Otevřít URL", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"d"}, KeyLabel: "d", Label: "Zvukové zařízení", Keymap: true},
	{Mode: commandModeMain, Keys: []string{"y"}, KeyLabel: "y", Label: "Text písně", Keymap: true},
	{Mode: commandModeMain | commandModeProvider | commandModeProviderPill | commandModeVolume | commandModeEQ | commandModeShuffle | commandModeRepeat | commandModeSpeed, Keys: []string{"tab", "shift+tab"}, KeyLabel: "Tab/Shift+Tab", Label: "Přepnout část", Keymap: true, ContextHelp: true},
	{Mode: commandModeMain, Keys: []string{"esc", "backspace", "b"}, KeyLabel: "Esc", Label: "Zpět na playlisty", Keymap: true, ContextHelp: true, Cancel: true},
	{Mode: commandModeAny, Keys: []string{"ctrl+k"}, KeyLabel: "Ctrl+K", Label: "Nápověda", Keymap: true, ContextHelp: true, Help: true},
	{Mode: commandModeMain, Keys: []string{"?"}, KeyLabel: "?", Label: "Nápověda", Keymap: true},
	{Mode: commandModeAny, Keys: []string{"ctrl+c"}, KeyLabel: "Ctrl+C", Label: "Konec", Keymap: true},
	// q queues inside the subscriptions overlay, so the keymap must not list
	// it as quit there.
	{Mode: commandModeAny, Keys: []string{"q"}, KeyLabel: "q", Label: "Konec", Keymap: true, Enabled: func(m Model) bool { return !m.subs.visible }},
	{Mode: commandModeAny, Keys: []string{"ctrl+z"}, KeyLabel: "Ctrl+Z", Label: "Vrátit změnu"},
	{Mode: commandModeProvider, Keys: []string{"ctrl+r"}, KeyLabel: "Ctrl+R", Label: "Obnovit", Keymap: true, ContextHelp: true},

	// Shared text editing is reserved even though these are intentionally absent
	// from the global keymap, where they would be misleading outside a field.
	{Mode: commandModeKeymapSearch | commandModeFileBrowserSearch | commandModeNavSearch | commandModePlaylistManagerInput | commandModePlaylistPickerInput | commandModeSearch | commandModeNetSearch | commandModeSpotSearch | commandModeJump | commandModeURL | commandModeThemePickerFilter | commandModeVisPickerFilter | commandModeProviderSearch, Keys: []string{"left", "right", "home", "end", "ctrl+a", "ctrl+e", "backspace", "delete", "ctrl+w", "ctrl+u"}, KeyLabel: "Text editor", Label: "Úpravy textu"},

	{Mode: commandModeProvider, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Načíst", LabelFor: func(m Model) string {
		if m.selectedProviderListIsBrowseEntry() {
			return "Otevřít"
		}
		return "Načíst"
	}, ContextHelp: true, Primary: true},
	{Mode: commandModeProvider, Keys: []string{"esc", "backspace", "b"}, KeyLabel: "Esc", Label: "Zpět", ContextHelp: true, Cancel: true},
	{Mode: commandModeProvider, Keys: []string{"l"}, KeyLabel: "l", Label: "Nejnovější díl do fronty", Keymap: true, ContextHelp: true, Enabled: func(m Model) bool {
		_, _, ok := m.selectedProviderShow()
		return ok
	}},
	{Mode: commandModeProvider, Keys: []string{"a"}, KeyLabel: "a", Label: "Přidat všechny díly", Keymap: true, Enabled: func(m Model) bool {
		_, _, ok := m.selectedProviderShow()
		return ok
	}},
	{Mode: commandModeProvider, Keys: []string{"f"}, KeyLabel: "f", Label: "Oblíbit", ContextHelp: true, Prominent: true, Enabled: func(m Model) bool {
		_, ok := m.provider.(provider.FavoriteToggler)
		if !ok || m.provLoading || m.provCursor < 0 || m.provCursor >= len(m.providerLists) || m.selectedProviderListIsBrowseEntry() {
			return false
		}
		if sl, ok := m.provider.(provider.SectionedList); ok {
			return sl.IsFavoritableID(m.providerLists[m.provCursor].ID)
		}
		return true
	}},
	{Mode: commandModeSpotSearch, Keys: []string{"f"}, KeyLabel: "f", Label: "Oblíbit", ContextHelp: true, Prominent: true, Enabled: func(m Model) bool {
		_, ok := m.spotSearch.prov.(provider.FavoriteToggler)
		return ok && m.spotSearch.screen == spotSearchResults && !m.spotSearchBusy() &&
			m.spotSearch.cursor >= 0 && m.spotSearch.cursor < len(m.spotSearch.results) &&
			m.spotSearch.results[m.spotSearch.cursor].IsAlbum() && m.spotSearch.results[m.spotSearch.cursor].AlbumID() != ""
	}},
	{Mode: commandModeEQ, Keys: []string{"up", "down"}, KeyLabel: "Up Down", Label: "Zesílení", ContextHelp: true},
	{Mode: commandModeSpeed, Keys: []string{"left", "right"}, KeyLabel: "Left Right", Label: "Rychlost", ContextHelp: true},
	{Mode: commandModeVolume, Keys: []string{"left", "right", "up", "down", "h", "l", "k", "j"}, KeyLabel: "Arrows", Label: "Hlasitost ±1 dB", ContextHelp: true, Primary: true},
	{Mode: commandModeShuffle, Keys: []string{"enter", "left", "right", "up", "down", "h", "l", "k", "j"}, KeyLabel: "Enter / Arrows", Label: "Náhodně", ContextHelp: true, Primary: true},
	{Mode: commandModeRepeat, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Opakování", ContextHelp: true, Primary: true},
	{Mode: commandModeRepeat, Keys: []string{"left", "right", "up", "down", "h", "l", "k", "j"}, KeyLabel: "Arrows", Label: "Režim opakování", ContextHelp: true},
	{Mode: commandModeProviderPill, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Otevřít", ContextHelp: true, Primary: true},
	{Mode: commandModeProviderPill, Keys: []string{"esc", "backspace"}, KeyLabel: "Esc", Label: "Zpět", ContextHelp: true, Cancel: true},
	{Mode: commandModeKeymap | commandModeFileBrowser | commandModeNavBrowser | commandModePlaylistManager | commandModePlaylistPicker | commandModeQueue | commandModeDevicePicker, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zpět", ContextHelp: true, Cancel: true},
	{Mode: commandModeKeymapSearch | commandModeFileBrowserSearch | commandModeNavSearch | commandModePlaylistManagerInput | commandModePlaylistPickerInput | commandModeSearch | commandModeNetSearch | commandModeSpotSearch | commandModeJump | commandModeURL | commandModeProviderSearch, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zrušit", ContextHelp: true, Cancel: true},
	{Mode: commandModeKeymapSearch | commandModeFileBrowserSearch | commandModeNavSearch | commandModePlaylistManagerInput | commandModePlaylistPickerInput | commandModeSearch | commandModeNetSearch | commandModeSpotSearch | commandModeJump | commandModeURL | commandModeProviderSearch, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Potvrdit", ContextHelp: true, Primary: true},
	{Mode: commandModeNavBrowser | commandModePlaylistManager | commandModePlaylistPicker | commandModeQueue | commandModeDevicePicker | commandModeProviderSearch, Keys: []string{"up", "down", "k", "j"}, KeyLabel: "Up Down", Label: "Pohyb", ContextHelp: true},
	{Mode: commandModeFileBrowser | commandModeNavBrowser | commandModePlaylistManager | commandModePlaylistPicker | commandModeDevicePicker, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Vybrat", ContextHelp: true, Primary: true},
	{Mode: commandModePlaylistPicker, Keys: []string{"p"}, KeyLabel: "p", Label: "Přidat na začátek", Keymap: true, ContextHelp: true},
	{Mode: commandModePlaylistManager, Keys: []string{"A"}, KeyLabel: "A", Label: "Přidat do aktuálního", Keymap: true, ContextHelp: true},
	{Mode: commandModeNavBrowser, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr", ContextHelp: true, Enabled: func(m Model) bool { return m.navBrowser.mode != navBrowseModeMenu }},
	{Mode: commandModeNavBrowser, Keys: []string{"f"}, KeyLabel: "f", Label: "Oblíbit", ContextHelp: true, Prominent: true, Enabled: func(m Model) bool {
		_, ok := m.navBrowser.prov.(provider.FavoriteToggler)
		idx := m.selectedNavRawIndex(len(m.navBrowser.albums))
		return ok && m.navView() == navViewAlbums && !m.navBrowser.loading && !m.navBrowser.albumLoading &&
			idx >= 0 && m.navBrowser.albums[idx].ID != ""
	}},
	{Mode: commandModeNavBrowser, Keys: []string{"f"}, KeyLabel: "f", Label: "Oblíbit žánr", LabelFor: func(m Model) string {
		if genre, ok := m.selectedNavGenre(); ok && genre.Favorite {
			return "Zrušit oblíbení žánru"
		}
		return "Oblíbit žánr"
	}, ContextHelp: true, Enabled: func(m Model) bool {
		_, canFavorite := m.navGenreBrowser().(provider.GenreFavoriteToggler)
		return canFavorite && m.navBrowser.mode == navBrowseModeByGenre && m.navBrowser.screen == navBrowseScreenList && !m.navBrowser.loading
	}},
	{Mode: commandModeKeymap, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr", ContextHelp: true, Primary: true},
	{Mode: commandModeThemePicker, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zrušit náhled", ContextHelp: true, Cancel: true},
	{Mode: commandModeThemePicker, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Použít vzhled", ContextHelp: true, Primary: true},
	{Mode: commandModeThemePicker, Keys: []string{"up", "down", "k", "j"}, KeyLabel: "Up Down", Label: "Náhled", ContextHelp: true},
	{Mode: commandModeThemePicker, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr", ContextHelp: true},
	{Mode: commandModeVisPicker, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zrušit náhled", ContextHelp: true, Cancel: true},
	{Mode: commandModeVisPicker, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Uložit", ContextHelp: true, Primary: true},
	{Mode: commandModeVisPicker, Keys: []string{"up", "down", "k", "j"}, KeyLabel: "Up Down", Label: "Náhled", ContextHelp: true},
	{Mode: commandModeVisPicker, Keys: []string{"/"}, KeyLabel: "/", Label: "Filtr", ContextHelp: true},
	{Mode: commandModeThemePickerFilter | commandModeVisPickerFilter, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zrušit filtr", ContextHelp: true, Cancel: true},
	{Mode: commandModeThemePickerFilter | commandModeVisPickerFilter, Keys: []string{"enter"}, KeyLabel: "Enter", Label: "Hotovo", ContextHelp: true, Primary: true},
	{Mode: commandModeQueue, Keys: []string{"d"}, KeyLabel: "d", Label: "Odebrat", Destructive: true, ContextHelp: true, Primary: true, Enabled: func(m Model) bool { return m.playlist != nil && m.playlist.QueueLen() > 0 }},
	{Mode: commandModeQueue, Keys: []string{"c"}, KeyLabel: "c", Label: "Vyprázdnit", Destructive: true, ContextHelp: true},
	{Mode: commandModeFileBrowser, Keys: []string{"R"}, KeyLabel: "R", Label: "Nahradit frontu", Destructive: true, ContextHelp: true},
	{Mode: commandModeNavBrowser, Keys: []string{"R"}, KeyLabel: "R", Label: "Nahradit frontu", Destructive: true, ContextHelp: true, Enabled: func(m Model) bool { return m.navView() == navViewTracks }},
	{Mode: commandModeLyrics, Keys: []string{"r"}, KeyLabel: "r", Label: "Znovu", ContextHelp: true, Primary: true, Enabled: func(m Model) bool { return !m.lyrics.loading && (m.lyrics.err != nil || len(m.lyrics.lines) == 0) }},
	{Mode: commandModeLyrics, Keys: []string{"[", "]"}, KeyLabel: "[ ]", Label: "Posun textu ±250 ms", ContextHelp: true, Keymap: true, Enabled: func(m Model) bool { return m.lyricsSyncable() && m.lyricsHaveTimestamps() }},
	{Mode: commandModeLyrics, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zavřít", ContextHelp: true, Cancel: true},
	{Mode: commandModeInfo, Keys: []string{"esc"}, KeyLabel: "Esc", Label: "Zavřít", ContextHelp: true, Cancel: true},
	{Mode: commandModePlaylistManager, Keys: []string{"a"}, KeyLabel: "a", Label: "Nový playlist", ContextHelp: true, Primary: true, Enabled: func(m Model) bool {
		return m.plManager.visible && m.plManager.screen == plMgrScreenList
	}},
	{Mode: commandModePlaylistManager, Keys: []string{"D"}, KeyLabel: "D", Label: "Přidat složky", ContextHelp: true, Enabled: func(m Model) bool {
		if !m.plManager.visible {
			return false
		}
		switch m.plManager.screen {
		case plMgrScreenTracks:
			return plMgrVirtualPlaylistName(m.plManager.selPlaylist) == ""
		case plMgrScreenList:
			idx := m.plMgrPlaylistRealIndex(m.plManager.cursor)
			return idx >= 0 && plMgrVirtualPlaylistName(m.plManager.playlists[idx].Name) == ""
		default:
			return false
		}
	}},
	{Mode: commandModePlaylistManager, Keys: []string{"f", "n"}, KeyLabel: "f/n", Label: "★/" + favHeart, ContextHelp: true, Enabled: func(m Model) bool {
		return m.plManager.visible && m.plManager.screen == plMgrScreenTracks
	}},
	{Mode: commandModePlaylistManager, Keys: []string{"[", "]"}, KeyLabel: "[ ]", Label: "Přeřadit", ContextHelp: true, Enabled: func(m Model) bool {
		return m.plManager.visible && m.plManager.screen == plMgrScreenTracks
	}},
	{Mode: commandModePlaylistManagerDirs, Keys: []string{"esc", "backspace", "h", "left"}, KeyLabel: "Esc", Label: "Zpět na skladby", ContextHelp: true, Cancel: true},
	{Mode: commandModePlaylistManagerDirs, Keys: []string{"a"}, KeyLabel: "a", Label: "Přidat složku", ContextHelp: true, Primary: true},
	{Mode: commandModePlaylistManagerDirs, Keys: []string{"d"}, KeyLabel: "d", Label: "Odebrat", Destructive: true, ContextHelp: true},
	{Mode: commandModePlaylistManagerDirs, Keys: []string{"r"}, KeyLabel: "r", Label: "Včetně podsložek", ContextHelp: true},
	{Mode: commandModePlaylistManagerDirs, Keys: []string{"up", "down", "k", "j"}, KeyLabel: "Up Down", Label: "Pohyb", ContextHelp: true},
	{Mode: commandModeFileBrowser, Keys: []string{"D"}, KeyLabel: "D", Label: "Přidat jako složku", ContextHelp: true, Enabled: func(m Model) bool {
		return m.fileBrowser.visible && m.fileBrowser.targetPlaylist != ""
	}},
}

func (m Model) commandHelp(mode commandMode) string {
	var cancel, primary, help, prominent, optional []commandSpec
	for _, command := range commandRegistry {
		if !command.ContextHelp || command.Mode&mode == 0 || !command.enabled(m) {
			continue
		}
		switch {
		case command.Cancel:
			cancel = append(cancel, command)
		case command.Primary:
			primary = append(primary, command)
		case command.Help:
			help = append(help, command)
		case command.Prominent:
			prominent = append(prominent, command)
		default:
			optional = append(optional, command)
		}
	}

	// Back/cancel, the primary action, and help are the only mandatory hints.
	// Keep them first so narrow terminals drop optional navigation rather than
	// trapping the user in an unfamiliar overlay.
	var ordered []commandSpec
	if len(cancel) > 0 {
		ordered = append(ordered, cancel[0])
	}
	if len(primary) > 0 {
		ordered = append(ordered, primary[0])
	}
	if len(help) > 0 {
		ordered = append(ordered, help[0])
	}
	ordered = append(ordered, prominent...)
	if (mode == commandModeMain || mode == commandModeProvider) && m.luaMgr != nil {
		// plugin keys (e.g. Ctrl+N YT radio) belong in the hint bar, not only in Ctrl+K
		for _, b := range m.luaMgr.KeyBindings() {
			if b.Description != "" {
				ordered = append(ordered, commandSpec{KeyLabel: pluginKeyLabel(b.Key), Label: b.Description})
			}
		}
	}
	ordered = append(ordered, optional...)
	return renderCommandHelp(ordered, m.helpWidth(), m)
}

func (m Model) helpWidth() int {
	if ui.PanelWidth > 0 {
		return ui.PanelWidth
	}
	if m.width > 0 {
		return m.width
	}
	return 80
}

func renderCommandHelp(commands []commandSpec, width int, m Model) string {
	if width <= 0 {
		return ""
	}
	narrow := width < 48
	var b strings.Builder
	for _, command := range commands {
		label := command.label(m)
		if narrow {
			switch {
			case command.Cancel:
				label = "Zpět"
			case command.Primary:
				label = "Spustit"
			case command.Help:
				label = "Nápověda"
			}
		}
		hint := helpKey(command.KeyLabel, label)
		if b.Len() > 0 {
			hint = " " + hint
		}
		if lipgloss.Width(b.String()+hint) > width {
			break
		}
		b.WriteString(hint)
	}
	return b.String()
}

// pluginKeyLabel turns a Bubbletea key string ("ctrl+n") into a hint label ("Ctrl+N").
func pluginKeyLabel(key string) string {
	parts := strings.Split(key, "+")
	for i, p := range parts {
		if len(p) == 1 {
			parts[i] = strings.ToUpper(p)
		} else if p != "" {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "+")
}
