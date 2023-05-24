package utils

// func pew(url string, size int64, data []byte) (http.Response, error) {

// 	uploadFile, err := os.Open("TODO.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer uploadFile.Close()
// 	info, err := uploadFile.Stat()
// 	if err != nil {
// 		panic(err)
// 	}
// 	buf := make([]byte, info.Size())
// 	_, err = bufio.NewReader(uploadFile).Read(buf)

// 	_, err = gay(request.URL, info.Size(), buf)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
// 	req.Header.Add("x-amz-acl", "public-read")
// 	if err != nil {
// 		// handle error
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		// handle error
// 	}
// 	defer resp.Body.Close()

// 	return *resp, nil
// 	// do something with the response
// }
