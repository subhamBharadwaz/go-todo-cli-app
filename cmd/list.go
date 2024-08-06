/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
	"todo-list/utils"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all todo items in the CSV file",
	Long: `The list command displays all todo items currently stored in the CSV file 'tasks.csv'. Each todo item is shown with its ID, description, and the time it was created. This command is useful for reviewing all the tasks you have added.

Example usage:
  $ tasks list

This command reads from the 'tasks.csv' file and formats the output in a tabular format for easy reading. If the file is empty or contains no todos, it will indicate that there are no tasks available. But this command will gives only the todos which done status is false.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the -all flag is set
		showAll, _ := cmd.Flags().GetBool("all")

		file, err := utils.LoadFile("tasks.csv")
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		defer utils.CloseFile(file)

		uncompletedTasks, err := utils.ReadTasks(file)
		if err != nil {
			fmt.Println("Error reading tasks", err)
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

		// Print headers based on --all flag
		if showAll {
			fmt.Fprintln(writer, "ID\tName\tCreated\tDone\t")
		} else {
			fmt.Fprintln(writer, "ID\tName\tCreated\t")
		}

		for _, task := range uncompletedTasks {
			// Print task details based on --all flag
			createdTime, _ := time.Parse(time.RFC3339, task[2])

			if showAll {
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n", task[0], task[1], timediff.TimeDiff(createdTime), task[3])

			} else {
				// check if the done is "false"
				if len(task) > 3 && task[3] == "false" {
					fmt.Fprintf(writer, "%s\t%s\t%s\t\n", task[0], task[1], timediff.TimeDiff(createdTime))
				}
			}

		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks, including completed ones")
}
