.PHONY: build-run build clean

build-run:
	make build
	./target/linux/valheim-syncer-server

build:
	go build -o target/linux/valheim-syncer-server main.go
	cp config/config.toml target/linux/

package-linux:
	make build
	cd target/linux && tar zcvf valheim-syncer-server-linux.tar.gz config.toml valheim-syncer-server && cd -

package-linux-installer:
	fyne package -os linux --release
	mkdir -p target/linux
	mv valheim-syncer-server.tar.xz target/linux/

package-windows:
	mkdir -p target/windows
	CC=x86_64-w64-mingw32-gcc fyne package -os windows --release --appID com.comoyi.valheim-syncer-server --name target/windows/valheim-syncer-server.exe
	cp config/config.toml target/windows/
	cd target/windows && zip valheim-syncer-server-windows.zip config.toml valheim-syncer-server.exe && cd -

clean:
	rm -rf target

bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go <font-file>
	#fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go ~/.local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf

bundle-font-build:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go /usr/local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf
