package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Strubbl/wallabago/v7"
	"github.com/spf13/cobra"
)

var (
	cfgUnread  bool
	cfgStarred bool
	cfgDate    string
	cfgDelete  bool
)

/*
  TODO: See comment for delete API call.
type wallabagURL struct {
	url string
}
*/

func init() {
	pruneCmd.PersistentFlags().BoolVarP(&cfgUnread, "unread", "u", false, "Include unread entries for deletion. False will prevent unread articles from being deleted")
	pruneCmd.Flags().BoolVarP(&cfgStarred, "starred", "s", false, "Include starred entry in deletion. False will prevent starred article to be deleted.")
	pruneCmd.Flags().StringVarP(&cfgDate, "date", "d", "", "Articles older than this date will be removed if they match the archived/starred flags, format \"YYYY-MM-DD\".")
	pruneCmd.Flags().BoolVar(&cfgDelete, "delete", false, "Delete articles. Without this flag, it will only do a dry run.")

	pruneCmd.MarkFlagRequired("date")

	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete old article from wallabag",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running prune command")

		home, _ := os.UserHomeDir()
		configJSON := home + "/.config/cleanABag/credentials.json"
		if cfgFile != "" {
			configJSON = cfgFile
		}
		err := wallabago.ReadConfig(configJSON)
		if err != nil {
			fmt.Println("Error reading credentials config file")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		baseDate, err := time.Parse("2006-01-02", cfgDate)
		if err != nil {
			fmt.Println("Wrong time format provided.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		nbEntries, err := wallabago.GetNumberOfTotalArticles()
		if err != nil {
			fmt.Println("Couldn't retrieve total number of articles.")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("There are", nbEntries, "saved on your wallabag instance.")

		url := wallabago.Config.WallabagURL +
			"/api/entries.json?perPage=" +
			strconv.Itoa(nbEntries) +
			"&detail=metadata" +
			"&sort=updated" +
			"&order=desc"
		response, err := wallabago.APICall(
			url,
			"GET",
			[]byte{},
		)
		if err != nil {
			fmt.Println("Couldn't retrieve articles from wallabag.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var e wallabago.Entries
		err = json.Unmarshal(response, &e)
		if err != nil {
			fmt.Println("Bad format response from wallabag.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		unread := 0
		if cfgUnread {
			unread = 1
		}
		starred := 0
		if cfgStarred {
			starred = 1
		}
		fmt.Println(
			"Will remove article older than",
			baseDate.Format("2006-01-02"),
			"and in status (unread, starred):",
			unread,
			starred,
		)
		var toRemove []wallabago.Item
		for i := 0; i < len(e.Embedded.Items); i++ {
			if e.Embedded.Items[i].UpdatedAt.Time.Before(baseDate) &&
				(e.Embedded.Items[i].IsArchived == 1 || unread == 1) &&
				(e.Embedded.Items[i].IsStarred == 0 || starred == 1) {
				toRemove = append(toRemove, e.Embedded.Items[i])
			}
		}

		if len(toRemove) == 0 {
			fmt.Println("Nothing to delete, leaving")
			os.Exit(0)
		}
		fmt.Println("This command will remove", len(toRemove), "entries:")
		//var urlToDelete []map[string]string
		for i := 0; i < len(toRemove); i++ {
			status := "  "
			if toRemove[i].IsArchived == 0 {
				status = "🆕"
			}
			if toRemove[i].IsStarred == 1 {
				status += "⭐"
			} else {
				status += "  "
			}

			fmt.Println(
				"- ",
				toRemove[i].UpdatedAt.Format("2006-01-02"),
				status,
				toRemove[i].Title,
			)
			//u := map[string]string{"url": toRemove[i].URL}
			//urlToDelete = append(urlToDelete, u)
		}

		if cfgDelete {
			// Seems there is an issue with the /entries/list.json
			// endpoint for the DELETE method
			// Instead, deleting one by one…
			// Keeping the following code to figure it out later.
			/*
				toDelete, _ := json.Marshal(urlToDelete)
				url := wallabago.Config.WallabagURL + "/api/entries/list.json"
				response, err := wallabago.APICall(
					url,
					"DELETE",
					toDelete,
				)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			*/
			baseURL := wallabago.Config.WallabagURL + "/api/entries/"
			for i := 0; i < len(toRemove); i++ {
				url := baseURL + strconv.Itoa(toRemove[i].ID)
				response, err := wallabago.APICall(
					url,
					"DELETE",
					[]byte{},
				)
				if err != nil {
					fmt.Println("Couldn't delete entry", toRemove[i].ID)
					fmt.Println(err.Error())
					os.Exit(1)
				}

				var item wallabago.Item
				err = json.Unmarshal(response, &item)
				if err != nil {
					fmt.Println("Bad format response from wallabag")
					fmt.Println(err.Error())
					os.Exit(1)
				}
				fmt.Println("Entry", item.Title, "(", item.URL, ") has been deleted.")
				time.Sleep(500 * time.Millisecond)
			}
		}
	},
}
