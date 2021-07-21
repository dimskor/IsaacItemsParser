package parser

import (
    "fmt"
    "github.com/microcosm-cc/bluemonday"
    "github.com/PuerkitoBio/goquery"
    "IsaacItemsParser/models"
    "log"
    "net/http"
    "net/url"
    "regexp"
    "strconv"
    "strings"
)

func getDoc(url string) *goquery.Document {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    if res.StatusCode != 200 {
        log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    return doc
}

func buildIconUrl(input string) string {
    imageUrl, err := url.Parse(input)
    if err != nil {
        log.Fatal(err)
    }

    path := imageUrl.Path
    path = strings.Split(path, "/revision")[0]
    path = strings.Trim(path, "/")
    path = strings.ReplaceAll(path, "?", "%3F")
    return fmt.Sprint("https://static.wikia.nocookie.net", "/", path)
}

func buildPageUrl(input string) string {
    pageUrl, err := url.Parse(input)
    if err != nil {
        log.Fatal(err)
    }

    path := pageUrl.Path
    path = strings.Trim(path, "/")
    return fmt.Sprint("https://bindingofisaacrebirth.fandom.com", "/", path)
}

func sanitizeDescription(input string) string {
    p := bluemonday.NewPolicy()
    p.AllowElements("br")
    p.AllowElements("img")
    p.AllowAttrs("alt").OnElements("img")
    p.AllowAttrs("title").OnElements("img")
    p.AllowAttrs("src").OnElements("img")
    input = p.Sanitize(input)

    re := regexp.MustCompile(`(/revision/latest/scale-to-width-down/[\d]+\?cb=[\d]+)`)
    input = re.ReplaceAllString(input, "")
    input = strings.TrimSpace(input)
    return input
}

func selectCollectibles(doc *goquery.Document) []models.TableRow {
    var rows []models.TableRow

    doc.Find(".wikitable tr.row-collectible").Each(func(i int, tr *goquery.Selection) {
        row := models.TableRow{}
        row.Name = tr.Find("td:nth-child(1)")
        row.Id = tr.Find("td:nth-child(2)")
        row.Icon = tr.Find("td:nth-child(3)")
        row.Quote = tr.Find("td:nth-child(4)")
        row.Description = tr.Find("td:nth-child(5)")
        row.Quality = tr.Find("td:nth-child(6)")
        row.Type = models.ItemTypeCollectible
        rows = append(rows, row)
    })

    return rows
}

func selectTrinkets(doc *goquery.Document) []models.TableRow {
    var rows []models.TableRow

    doc.Find(".wikitable tr.row-trinket").Each(func(i int, tr *goquery.Selection) {
        row := models.TableRow{}
        row.Name = tr.Find("td:nth-child(1)")
        row.Id = tr.Find("td:nth-child(2)")
        row.Icon = tr.Find("td:nth-child(3)")
        row.Quote = tr.Find("td:nth-last-child(2)")
        row.Description = tr.Find("td:last-child")
        row.Type = models.ItemTypeTrinket
        rows = append(rows, row)
    })

    return rows
}

func selectCardsRunes(doc *goquery.Document) []models.TableRow {
    var rows []models.TableRow

    doc.Find(".wikitable tr.row-pickup").Each(func(i int, tr *goquery.Selection) {
        row := models.TableRow{}

        row.Name = tr.Find("td:nth-child(1)")
        row.Id = tr.Find("td:nth-child(2)")
        row.Icon = tr.Find("td:nth-child(3)")
        row.Quote = tr.Find("td:nth-last-child(2)")
        row.Description = tr.Find("td:last-child")

        name := row.Name.Text()
        id := strings.TrimSpace(row.Id.Text())

        if (strings.Contains(name, "Rune")) {
            row.Type = models.ItemTypeRune
        } else if (strings.Contains(name, "Soul")) {
            row.Type = models.ItemTypeSoulstone
        } else if id == "5.300.49" || id == "5.300.50" || id == "5.300.78" {
            row.Type = models.ItemTypeOther
        } else {
            row.Type = models.ItemTypeCard
        }

        rows = append(rows, row)
    })

    return rows
}

func parseTable(rows []models.TableRow) []models.Item {
    items := []models.Item{}

    for _, row := range rows {
        item := models.Item{}

        // Id
        item.Id = strings.TrimSpace(row.Id.Text())

        // IconUrl
        iconUrl, ok := row.Icon.Find("a > img").Attr("src")
        if ok == false {
            log.Fatal("Selector error: IconUrl")
        }
        item.IconUrl = buildIconUrl(iconUrl)

        // Name
        item.Name = strings.TrimSpace(row.Name.Text())

        // PageUrl
        pageUrl, ok := row.Name.Find("a").Attr("href")
        if ok == false {
            if row.Type == models.ItemTypeCard || row.Type == models.ItemTypeRune || row.Type == models.ItemTypeSoulstone {
                item.PageUrl = "https://bindingofisaacrebirth.fandom.com/wiki/Cards_and_Runes"
            } else {
                log.Fatal("Selector error: PageUrl")
            }
        } else {
            item.PageUrl = buildPageUrl(pageUrl)
        }

        // Charges, Type
        charges, ok := row.Icon.Find("div:last-child > img").Attr("title")
        if ok == false {
            if row.Type == models.ItemTypeCollectible {
                item.Type = models.ItemTypePassive
            } else {
                item.Type = row.Type
            }
        } else {
            item.Type = models.ItemTypeActive
            charges = strings.ToLower(charges)
            item.Charges = &charges
        }

        // Quote
        item.Quote = strings.TrimSpace(row.Quote.Text())

        // Description
        description, err := row.Description.Html()
        if err != nil {
            log.Fatal("Selector error: Description")
        }
        item.Description = sanitizeDescription(description)

        // Quality
        if row.Quality != nil {
            quality := strings.TrimSpace(row.Quality.Text())
            if quality != "" {
                qualityInt, err := strconv.Atoi(quality)
                if err != nil {
                    log.Fatal(err)
                }
                item.Quality = &qualityInt
            }
        }

        items = append(items, item)
    }

    return items
}

func Parse() []models.Item {
    var items []models.Item

    // Collectibles
    collectiblesDoc := getDoc("https://bindingofisaacrebirth.fandom.com/wiki/Items")
    collectiblesTable := SelectCollectibles(collectiblesDoc)
    collectibles := ParseTable(collectiblesTable)
    items = append(items, collectibles...)
    fmt.Printf("Parsed collectibles: %d\n", len(collectibles))

    // Trinkets
    trinketsDoc := getDoc("https://bindingofisaacrebirth.fandom.com/wiki/Trinkets")
    trinketsTable := SelectTrinkets(trinketsDoc)
    trinkets := ParseTable(trinketsTable)
    items = append(items, trinkets...)
    fmt.Printf("Parsed trinkets: %d\n", len(trinkets))

    // Cards, Runes
    cardsRunesDoc := getDoc("https://bindingofisaacrebirth.fandom.com/wiki/Cards_and_Runes")
    cardsRunesTable := SelectCardsRunes(cardsRunesDoc)
    cardsRunes := ParseTable(cardsRunesTable)
    items = append(items, cardsRunes...)
    fmt.Printf("Parsed cards & runes: %d\n", len(cardsRunes))

    return items
}
