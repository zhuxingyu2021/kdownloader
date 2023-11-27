package downloader

import (
	"context"
	"kdownloader/pkg/downloader"
	"testing"
	"time"
)

func TestDownloadFile(t *testing.T) {
	url := `https://c4.kemono.su/data/30/1b/301bbf6f37da223270e7070eb6a45219e07cef8c5fb50f873b7cccff694dda9e.jpg?f=Untitled_1.jpe`
	path := `1.jpe`

	ctx := downloader.DownloadContext(context.Background())
	err := downloader.DownloadNormalFile(ctx, url, path)

	if err != nil {
		panic(err)
	}
}

func TestDownloadWorker(t *testing.T) {
	urls := []string{
		"https://c1.kemono.su/data/64/b4/64b4a8e823cc1c2912efc5932d345cd2239aafb7507f8bd7b64e7f5c58d969db.jpg?f=04550.jpg",
		"https://c1.kemono.su/data/64/b4/64b4a8e823cc1c2912efc5932d345cd2239aafb7507f8bd7b64e7f5c58d969db.jpg?f=04550.jpg",
		"https://c1.kemono.su/data/5d/6d/5d6d269d54740424fae658edad36534bcb7b19fd76b0eefd07a2c3437bdcb26a.jpg?f=04546.jpg",
		"https://c4.kemono.su/data/76/b1/76b14ade168661c9472af2cc212f808f46fdcdc3f00096b43ac0ea272d3a84d0.jpg?f=04547.jpg",
		"https://c4.kemono.su/data/75/f3/75f338a230377f3c4f3d05d13e37c52cda1cb0757b6b6535fdd08ce98842cb81.jpg?f=04548.jpg",
		"https://c3.kemono.su/data/0a/c2/0ac24deb2a317a204bd2328f9ec6349b4254331bf0e4166e8d67d6f1374a5077.jpg?f=04549.jpg",
		"https://c4.kemono.su/data/c9/13/c913cd2f1eefd93d766116c34d9e38045d84b3d7b496ca99dffb0b2f64cbceec.jpg?f=04550.jpg",
		"https://c3.kemono.su/data/bf/df/bfdf90dda69d4120212eaf8e583d1e94978797bfd54bd94d80a1795962fa1eab.jpg?f=04551.jpg",
		"https://c2.kemono.su/data/2f/a6/2fa60c59f2806cbbf111011ae44ec005a893b8ec86126ba311f318421ed5f846.jpg?f=04552.jpg",
		"https://c4.kemono.su/data/72/62/7262034848211b66f5a59f2489f576f1bc11697379d6cd542fe661bda289da4c.jpg?f=04553.jpg",
		"https://c2.kemono.su/data/a2/99/a29900b34fc40bd28214935e962e2aeba9b8b8c80a41c84cf11abe9138d7cfdb.jpg?f=04554.jpg",
		"https://c6.kemono.su/data/ac/20/ac201b2e5dee82c3047ec37fa8f6a68fc67fcbe00e50b713ab4bbfce421a56cf.jpg?f=04555.jpg",
		"https://c6.kemono.su/data/ed/0b/ed0b49f54c5cde0b823ec101af4d7ebd483618b1467aa33cf7affcda7463ab3c.jpg?f=04556.jpg",
		"https://c2.kemono.su/data/7b/54/7b54a0003b67b7ab33a13024891f67dee7f5b243989f82e00a10806d78afd6d9.jpg?f=04557.jpg",
		"https://c2.kemono.su/data/1b/2a/1b2a685bb97b6fe48d41243e868f0c870a2726e4d28bc0722c4651effd472709.jpg?f=04558.jpg",
		"https://c2.kemono.su/data/dc/3d/dc3dc0fe6bae7b618f0bb5ca46df8c25288ad7bd762baf786101688d269c7bb7.jpg?f=04559.jpg",
		"https://c1.kemono.su/data/78/eb/78eb6588aaf1648df8f2a5113e5f5e5ba5f63d855ab60d4cd8b6cc3b9b51913c.jpg?f=04560.jpg",
		"https://c4.kemono.su/data/0a/a0/0aa03fcb7251f08f34de8d1477adeb0a83851d2241884c1864febdec2636b9b7.jpg?f=04561.jpg",
		"https://c1.kemono.su/data/03/77/0377dd004af503eff4d2d45acf24b9013c4709779caf963a5877e359388842bb.jpg?f=04562.jpg",
		"https://c6.kemono.su/data/b9/cb/b9cb2c9b12cd6bc10bb2fc332b03496b934955a4f90efb955712c22cb0149e2a.jpg?f=04563.jpg",
		"https://c1.kemono.su/data/ff/51/ff51b01bb732575a15a4d2138b2afd59acc81a23ee548b74f95ded4b6d65d24f.jpg?f=04564.jpg",
		"https://c6.kemono.su/data/77/5b/775bdf8464bdd3746c4e3689a202aa7bf4e76a01ccac664d04411bf7b7f67360.jpg?f=04565.jpg",
		"https://c2.kemono.su/data/d0/c2/d0c2eedec1529a970661e3cc787fe3bd14c032d930a6d3e179e666a0aa834a5d.jpg?f=04566.jpg",
		"https://c5.kemono.su/data/5c/69/5c69e5e45f33604d6313a5218c3e64c4c5ef0ca1814cd73955a5862309bb107b.jpg?f=04567.jpg",
		"https://c6.kemono.su/data/35/70/357079bd6f2499d50c39b842475e87da5b4483a019332c45bb486521e835eb08.jpg?f=04568.jpg",
		"https://c5.kemono.su/data/06/af/06afead18535c972760b4a4dc350ad080a8dcafb0957ce32be531283be8ec090.jpg?f=04569.jpg",
		"https://c6.kemono.su/data/44/9f/449feaa3654c46571fd82d5c139431b9881c0966249849aecfd27229908254d2.jpg?f=04570.jpg",
		"https://c3.kemono.su/data/2d/c6/2dc6e6bf035f7081f9b3b3254a9d053171ebbf34ac41a9114738aff2eb7682a5.jpg?f=04571.jpg",
		"https://c4.kemono.su/data/d2/7f/d27f93bb71d78dfad026cbd308b5186a7bd8304c9a66bd560d0e0da6f57c42e1.jpg?f=04572.jpg",
		"https://c5.kemono.su/data/2f/95/2f95ef57f5e959dbc4303196bd52767bc67aa7def35d0704d0a3c4924c7b709b.jpg?f=04573.jpg",
		"https://c1.kemono.su/data/69/0a/690a6108dc6ec072908e93988c1db8dcb928958dbcb29cf4753e72c883d79823.jpg?f=04574.jpg",
		"https://c1.kemono.su/data/62/ed/62ed866541f4025781863fc8c5cd76fb3d854a2cf043e7ce15ae88cf30bb8515.jpg?f=04575.jpg",
		"https://c2.kemono.su/data/38/2d/382d6c9ba58c13b9c4e99cf33e4dccc69e07aa7fcde78eb3cfed48f21b7a352d.jpg?f=04576.jpg",
		"https://c4.kemono.su/data/09/b6/09b67c50f3fad4e87b6fb4dade12737cc2f471604d147f21df40d523fabcf8fa.jpg?f=04577.jpg",
		"https://c1.kemono.su/data/c3/56/c35645419f11686e83dc49a81a7e5ec379fc1638c4bb3678597dcbb87f13318f.jpg?f=04578.jpg",
		"https://c6.kemono.su/data/38/9c/389c8482ec5421abe095d2dfdf3bd3196133b5cdd56daf00e54f380f89b58281.jpg?f=04579.jpg",
		"https://c1.kemono.su/data/92/bf/92bf381d42883f8ca55e9e913b9e7d63d78aa55cdad36788f53296c21a769105.jpg?f=04580.jpg",
		"https://c4.kemono.su/data/cf/e9/cfe9ef98ba4d61cb43a059ce333f0a69e001606d4a42de7fca0f540d79992356.jpg?f=04581.jpg",
		"https://c5.kemono.su/data/20/6d/206d799e2fab84c0441ca929ea6057c0370dff7d21ca4734b63f6ede69a1900c.jpg?f=04582.jpg",
		"https://c6.kemono.su/data/a5/59/a5595213047b6ec6cfe541c2958c9b90ffdf3cc7e239e8b5dace5eebca33de02.jpg?f=04583.jpg",
		"https://c1.kemono.su/data/e2/88/e288967d49b9293880beec671da4dbcec82c245d0f6c5d3646bd4c42616b832c.jpg?f=04584.jpg",
		"https://c4.kemono.su/data/8a/1e/8a1e711afcbc8fa77b43bf72cecb3d7c700b8ff33460ed8a07d9fbf6ee41822b.jpg?f=04585.jpg",
	}

	ctx_, cancel := context.WithCancel(context.Background())
	ctx := downloader.DownloadContext(ctx_)

	urlchan := make(chan (string))

	done := make(chan (bool))
	go func() {
		downloader.DWorker(ctx, urlchan)
		done <- true
	}()

	for _, url := range urls {
		urlchan <- url
	}

	time.Sleep(time.Second)
	for len(downloader.GetUndownloadUrls(ctx)) > 0 {
		time.Sleep(time.Second)
	}

	cancel()
	<-done

	files := downloader.ListOKUrls(ctx)

	println("urls: ", len(urls), "downloads: ", len(files))

}
