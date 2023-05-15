package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

func find_down_pixels_down(data [][]string, line int , pos int) (int){
	for i := line; i < len(data) && i < line + 7; i++ {
		if(data[i][pos] == "1"){
			return i;
		}
	}
	return -1;
}

func find_down_pixels_up(data [][]string, line int , pos int) (int){
	for i := line; i > 0 && i > line - 7; i-- {
		if(data[i][pos] == "1"){
			return i;
		}
	}
	return -1;
}

func checkLineForEmptyNextPixels(line []string, pos int) bool{
	for i := pos; i < pos + 20; i++ {
		if i >= len(line) {
			return false;
		}
		if line[i] == "0"{
			return true;
		}
	}
	return false;
}

func checkForPixelsToBeFilledVertical(tmp[][]string, x int, y int) bool {

	if y == 0 {
		return true;
	}

	for i := y; i > 0 && i > y - 10; i-- {
		if tmp[i][x] == "1" {
			return false;
		}
	}
	return true;
}

func checkForPixelsToBeFilledHorizontal(tmp []string, x int) bool {

	if x == 0 {
		return true;
	}

	for i := x; i > 0 && i > x - 10; i-- {
		if tmp[i] == "1" {
			return false;
		}
	}
	return true;

}

func fillGapsInVerticalWalls(tmp [][]string, x int, y int) (bool, int) {
	for i := y + 2 ; i < len(tmp) && i <= y + 15; i++ {
		if tmp[i][x] == "1" {
			return true, i;
		}
	}
	return false, -1;
}

func writeImage(tmp [][]string, filename string){
	
	upLeft := image.Point{0, 0}
	lowRight := image.Point{len(tmp[0]), len(tmp)}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	white := color.RGBA{255, 255, 255, 0xff}
	black := color.RGBA{0, 0, 0, 0xff}
	
	for i := 0; i < len(tmp); i++ {
		for j := 0; j < len(tmp[0]); j++ {
			if tmp[i][j] == "1" {
				img.Set(j, i, white)
			}else if tmp[i][j] == "0" {
				img.Set(j, i, black);
			}
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img);
}

func main() {
	inputFile := "map.txt"
 
	filebuffer, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inputdata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)
 
	var mapData [][]string;
	var tmp []string;
	for data.Scan() {
		if data.Text() == "\n" {
			mapData = append(mapData, tmp);			
			tmp = []string{}
		}else {
			tmp = append(tmp, data.Text())
		}
	}

	var new_map_data [][]string;

	for i := 0; i < len(mapData); i++ {
		tmp = []string{}
		for j := 0; j < len(mapData[0]); j++ {
			if mapData[i][j] == "1" && checkForPixelsToBeFilledVertical(new_map_data, j , i - 1){
				tmp = append(tmp, "1");
			}else {
				tmp = append(tmp, "0");
			}
		}
		new_map_data = append(new_map_data, tmp);
	}
	writeImage(new_map_data, "map-wavy.png");


	// fix waving pixels 
	for i := 0; i < len(new_map_data); i++{
		for j := 0; j < len(new_map_data[0]); j++ {
			if index := find_down_pixels_down(new_map_data, i, j); j > 0 && index != -1 && new_map_data[i][j - 1] == "1"  {
				new_map_data[index][j] = "0";
				new_map_data[i][j] = "1";
			}
			if index := find_down_pixels_up(new_map_data, i, j); j > 0 && index != -1 && new_map_data[i][j - 1] == "1"  {
				new_map_data[index][j] = "0";
				new_map_data[i][j] = "1";
			}
		}
	}
	writeImage(new_map_data, "map-post-wavy.png");


	// init 
	final_map_data := make([][]string, len(mapData));
	for i := 0; i < len(new_map_data); i++ {
		final_map_data[i] = make([]string, len(mapData[0]))
	}

	// fillup 
	z := 0;
	for i := 0; i < len(new_map_data[0]); i++ {
		for j := 0; j < len(new_map_data); j++ {
			if z > 0 {
				final_map_data[j][i] = "1";
				/*
				if final_map_data[j][i] != "0" {
					tmp_index := i + 1;
					for tmp_index < i + 10 && tmp_index < len(new_map_data[0]) {
						final_map_data[j][tmp_index] = "0";
						new_map_data[j][tmp_index] = "0";
						tmp_index++;
					}
				}*/
				z--;
			}else if isWritable, yindex:= fillGapsInVerticalWalls(new_map_data, i, j); new_map_data[j][i] == "1" && isWritable {
				z = 10;
				// replace added  horizontal line with spaces 
				if yindex != -1 {
					tmp_index := i + 1;
					for new_map_data[yindex][tmp_index] == "1" && tmp_index < len(new_map_data[0]) {
						final_map_data[yindex][tmp_index] = "0";
						new_map_data[yindex][tmp_index] = "0";
						tmp_index++;
					} 
				}
				final_map_data[j][i] = "1";
			}else if new_map_data[j][i] == "1" {
				final_map_data[j][i] = "1"
			}else {
				final_map_data[j][i] = "0";
			}
		}
		z = 0;
	} 

	/*
	new_map_data = [][]string{}
	for i := 0; i < len(mapData); i++ {
		tmp = []string{}
		for j := 0; j < len(mapData[0]); j++ {
			if mapData[i][j] == "1" && checkForPixelsToBeFilledHorizontal(tmp, j - 1){
				tmp = append(tmp, "1");
			}else {
				tmp = append(tmp, "0");
			}
		}
		new_map_data = append(new_map_data, tmp);
	}
	*/

	writeImage(final_map_data, "map-final.png");
}
