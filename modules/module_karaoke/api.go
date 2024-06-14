package module_karaoke

import (
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
)

type ApiResponse struct {
	Meta struct {
		Page    int `json:"page"`
		MaxPage int `json:"max_page"`
		Total   int `json:"total"`
	} `json:"meta"`
	Results []PhkSong `json:"results"`
}

type SongResult struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
	Cover  string `json:"cover"`
}

var VALID_FORMATS = []string{"CDG", "WEBM", "MP4"}
var VALID_COVER_TYPE = []string{"NO_COVER", "UPLOADED", "LINK"}
var nonAsciiRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// @TODO: Rework temporarly to make this work
// This will later be dropped when we have a separate API to do this
func songPost(c echo.Context) error {
	/*
		//#region Parsing form
		formData := new(DtoSongCreate)
		if err := c.Bind(formData); err != nil {
			return err
		}

		if err := c.Validate(formData); err != nil {
			return err
		}
		//#endregion

		//#region Building folder name
		t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		titleFoldername, _, _ := transform.String(t, formData.Title)
		artistFoldername, _, _ := transform.String(t, formData.Artist)

		titleFoldername = strings.ReplaceAll(titleFoldername, " ", "")
		artistFoldername = strings.ReplaceAll(artistFoldername, " ", "")

		titleFoldername = nonAsciiRegex.ReplaceAllString(titleFoldername, "")
		artistFoldername = nonAsciiRegex.ReplaceAllString(artistFoldername, "")

		foldername := strings.ToLower(artistFoldername + "_" + titleFoldername)
		//#endregion

		// Need to test if validation is working correctly
		songField, err := c.FormFile("song")
		if err != nil || songField == nil {
			fmt.Println("Err: ", err)
			// Should not happen but meh
			return echo.NewHTTPError(http.StatusBadRequest, "Song missing!")
		}

		cdgField, _ := c.FormFile("cdg")
		if formData.Format == "CDG" && cdgField == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing CDG file")
		}

		//#region Creating the song in DB
		dbSong, err := ormCreateSong(
			uuid.New().String(),
			foldername,
			formData.Artist,
			formData.Title,
			strings.ToLower(formData.Format),
			"", // @TODO
		)

		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create song: "+err.Error())
		}
		//#endregion

		tempDir, err := os.MkdirTemp("", "phkaraoke")
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create temp dir: "+err.Error())
		}

		//#region Uploading the song
		songFilename := "song." + strings.ToLower(formData.Format)
		if formData.Format == "CDG" {
			songFilename = "song.mp3"
		}

		song, err := songField.Open()
		if err != nil {
			return err
		}
		defer song.Close()

		outputSong, err := os.Create(filepath.Join(tempDir, songFilename))
		if err != nil {
			song.Close()
			return err
		}

		if _, err := io.Copy(outputSong, song); err != nil {
			song.Close()
			outputSong.Close()
			return err
		}

		song.Close()
		outputSong.Close()
		//#endregion

		//#region Uploading the CDG when present
		if formData.Format == "CDG" {
			cdg, err := cdgField.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open CDG: "+err.Error())
			}

			outputCdg, err := os.Create(filepath.Join(tempDir, "song.cdg"))
			if err != nil {
				cdg.Close()
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create temp CDG file: "+err.Error())
			}

			_, err = io.Copy(outputCdg, cdg)
			if err != nil {
				cdg.Close()
				outputCdg.Close()

				fmt.Println(err)
				return err
			}

			cdg.Close()
			outputCdg.Close()
		}
		//#endregion

		coverPath := filepath.Join(tempDir, "cover.jpg")
		outputCover, err := os.Create(coverPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create temp cover file: "+err.Error())
		}

		//#region Uploading cover
		if formData.CoverType == "UPLOADED" {
			coverField, err := c.FormFile("cover")
			if coverField == nil || err != nil {
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get cover file: "+err.Error())
			}

			cover, err := coverField.Open()
			if err != nil {
				outputCover.Close()
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open cover file: "+err.Error())
			}

			_, err = io.Copy(outputCover, cover)
			if err != nil {
				outputCover.Close()
				cover.Close()
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy cover file: "+err.Error())
			}

			cover.Close()
		}
		//#endregion

		//#region Getting the cover from a URL
		if formData.CoverType == "LINK" {
			resp, err := http.Get(*formData.CoverUrl)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to download cover URL: "+err.Error())
			}

			_, err = io.Copy(outputCover, resp.Body)
			if err != nil {
				resp.Body.Close()
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy cover downloaded: "+err.Error())
			}
			resp.Body.Close()
		}
		//#endregion

		outputCover.Close()

		//#region Converting & Resizing cover
		if formData.CoverType != "NO_COVER" {
			buf, err := bimg.Read(coverPath)
			if err != nil {
				fmt.Println("Failed to open image: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read cover: "+err.Error())
			}

			newImage, err := bimg.NewImage(buf).Resize(300, 300)
			if err != nil {
				fmt.Println("Failed to resize image: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to resize cover: "+err.Error())
			}

			format := bimg.NewImage(newImage).Type()
			if format != "jpeg" {
				fmt.Printf("Wrong format: %v, expected %v. Converting...\n", format, "jpeg")
				newImage, err = bimg.NewImage(newImage).Convert(bimg.JPEG)
				if err != nil {
					fmt.Println("Failed to convert image: ", err)
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert cover: "+err.Error())
				}
			}

			err = bimg.Write(coverPath, newImage)
			if err != nil {
				fmt.Println("Failed to save resized image: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save resized cover: "+err.Error())
			}
		}
		//#endregion

		//#region Adding info.txt file
		infoFile, err := os.Create(filepath.Join(tempDir, "info.txt"))
		if err != nil {
			fmt.Println("Failed to create info.txt file: ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create info.txt file: "+err.Error())
		}
		infoFile.WriteString(formData.Artist + "\n" + formData.Title + "\n" + strings.ToLower(formData.Format) + "\n0")
		infoFile.Close()
		//#endregion

		err = utils.CopyDir(tempDir, filepath.Join(config.GET.RootPath, "karaoke", foldername))
		if err != nil {
			fmt.Println("Failed to copy song to the main directory: ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy song folder: "+err.Error())
		}

		os.RemoveAll(tempDir)

		return c.JSON(
			http.StatusCreated,
			songGet(dbSong),
		)
	*/
	return c.String(http.StatusInternalServerError, "Not implemented yet")
}

func spotifySearch(c echo.Context) error {
	query := c.QueryParam("q")
	if len(query) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	tracks, err := services.GET.Spotify.SearchSong(query)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to search on spotify: "+err.Error(),
		)
	}

	songs := []SongResult{}
	for _, t := range tracks {
		var bestImageUrl string = ""
		image := getBestImage(t.Album.Images)

		if image != nil {
			bestImageUrl = image.URL
		}

		song := SongResult{
			Artist: t.Artists[0].Name,
			Song:   t.Name,
			Cover:  bestImageUrl,
		}

		songs = append(songs, song)
	}

	return c.JSON(
		http.StatusOK,
		songs,
	)
}

func searchSongs(c echo.Context) error {
	query := c.QueryParam("q")
	pageStr := c.QueryParam("page")
	var page int = 1

	if len(pageStr) > 0 {
		var err error = nil
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}

	songs, err := ormListSongs(
		query,
		(page-1)*CONFIG.AmtSongsPerPage,
		CONFIG.AmtSongsPerPage,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"err":     "Failed to list songs",
			"details": err.Error(),
		})
	}

	count, err := ormCountSongs(query)
	if err != nil {
		count = len(songs)
	}

	return c.JSON(http.StatusOK, models.ContextualizedResponse{
		Results: songs,
		Meta: models.ResponseMetadata{
			LastPage: int(math.Ceil(float64(count) / float64(CONFIG.AmtSongsPerPage))),
			Total:    count,
		},
	})
}

func rescanSongs(c echo.Context) error {
	if err := INSTANCE.ScanSongs(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
