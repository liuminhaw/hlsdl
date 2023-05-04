/*
Copyright © 2023 Min-Haw, Liu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"

	"github.com/liuminhaw/hlsdl"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "downloading m3u8 file and all its ts segments into a single TS file",
	RunE:  cmdF,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("download called")
	// 	cmfF
	// },
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	downloadCmd.Flags().StringP("url", "u", "", "The manifest (m3u8) url")
	downloadCmd.Flags().StringP("dir", "d", "./download", "The directory where the file will be stored")
	downloadCmd.Flags().BoolP("record", "r", false, "Indicate whether the m3u8 is a live stream video and you want to record it")
	downloadCmd.Flags().IntP("workers", "w", 2, "Number of workers to execute concurrent operations")
	downloadCmd.SetArgs(os.Args[1:])
}

func cmdF(command *cobra.Command, args []string) error {
	m3u8URL, err := command.Flags().GetString("url")
	if err != nil {
		return err
	}

	dir, err := command.Flags().GetString("dir")
	if err != nil {
		return err
	}

	workers, err := command.Flags().GetInt("workers")
	if err != nil {
		return err
	}

	if record, err := command.Flags().GetBool("record"); err != nil {
		return err
	} else if record {
		return recordLiveStream(m3u8URL, dir)
	}

	return downloadVodMovie(m3u8URL, dir, workers)
}

func downloadVodMovie(url string, dir string, workers int) error {
	hlsDL := hlsdl.New(url, nil, dir, workers, true, "")
	filepath, err := hlsDL.Download()
	if err != nil {
		return err
	}
	log.Println("Downloaded file to " + filepath)
	return nil
}

func recordLiveStream(url string, dir string) error {
	recorder := hlsdl.NewRecorder(url, dir)
	recordedFile, err := recorder.Start()
	if err != nil {
		os.RemoveAll(recordedFile)
		return err
	}

	log.Println("Recorded file at ", recordedFile)
	return nil
}
