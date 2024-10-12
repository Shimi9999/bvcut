# ffmpegビルド手順メモ
ffmpeg内蔵のAACエンコーダは音質が悪いので、音質の良いfdk_aac入りのffmpegを自分でビルドする。

## 手順
基本的に以下のページをそのまま参考にして進めていけばOK。  
https://blog.tsukumijima.net/article/ffmpeg-windows-build/

## 環境
Windows11

## つまづいたポイント

### WSL (Windows Subsystem for Linux) を入れる - Ubuntuインストール時

#### エラー1
```
Installing, this may take a few minutes… WslRegisterDistribution failed with 
error: 0x800701bc Error: 0x800701bc WSL 2 ???????????? ??????????????????????? https://aka.ms/wsl2kernel ?????????
```

**解決法**  
Linux カーネル更新プログラム パッケージをダウンロードする  
https://learn.microsoft.com/ja-jp/windows/wsl/install-manual#step-4---download-the-linux-kernel-update-package

#### エラー2
```
WslRegisterDistribution failed with error: 0x80370102
Please enable the Virtual Machine Platform Windows feature and ensure virtualization is enabled in the BIOS.
For information please visit https://aka.ms/enablevirtualization
```

**解決法**  
1. 「Windowsの機能の有効化または無効化」を開く。
2. 「仮想マシン プラットフォーム」にチェックをつけた後、Windowsを再起動する。

参考: https://note.com/calm_hornet363/n/n7be9a247dcc0

### ビルドする - ビルドコマンド実行時
自分が最終的に使用したビルドコマンド
```
./cross_compile_ffmpeg.sh --build-ffmpeg-static=y --disable-nonfree=n --ffmpeg-git-checkout-version=n4.4.4
```
ビルドには時間がかかる

#### エラー1
```
Could not find the following execs (svn is actually package subversion, makeinfo is actually package texinfo if you're missing them): libtoolize cmake python clang meson bzip2 autogen gperf nasm unzip pax g++ makeinfo bison flex cvs yasm automake autoconf gcc svn make pkg-config ragel
Install the missing packages before running this script.
for ubuntu:
$ sudo apt-get update
$ sudo apt-get install subversion ragel curl texinfo g++ ed bison flex cvs yasm automake libtool autoconf gcc cmake git make pkg-config zlib1g-dev unzip pax nasm gperf autogen bzip2 autoconf-archive p7zip-full meson clang python3-distutils python-is-python3 -y
NB if you use WSL Ubuntu 20.04 you need to do an extra step: https://github.com/rdp/ffmpeg-windows-build-helpers/issues/452
```

**解決法**  
必要なexecsを取得するために、`$ sudo apt-get install`の行のコマンドをそのまま実行する。
```
sudo apt-get install subversion ragel curl texinfo g++ ed bison flex cvs yasm automake libtool autoconf gcc cmake git make pkg-config zlib1g-dev unzip pax nasm gperf autogen bzip2 autoconf-archive p7zip-full meson clang python3-distutils python-is-python3 -y
```

#### エラー2
最新版にした方がいいと思ってffmpegのバージョンに6.x.xを指定するとエラーが起きる。
```
./cross_compile_ffmpeg.sh --build-ffmpeg-static=y --disable-nonfree=n --ffmpeg-git-checkout-version=n6.1.1
```
エラーメッセージ
```
...
checking size of int *...
configure: error: Xvid supports only 32/64 bit architectures
failed configure generic
```
ffmpeg-windows-build-helpersが4.x.xしか対応していないっぽい?

**解決法**  
バージョン4.x.xでビルドする。

#### エラー3
```
libavcodec/libsvtav1.c: In function 'alloc_buffer':
libavcodec/libsvtav1.c:124:51: error: 'EbSvtAv1EncConfiguration' has no member named 'compressed_ten_bit_format'
  124 |         (config->encoder_bit_depth > 8) && (config->compressed_ten_bit_format == 0) ? 1 : 0;
      |                                                   ^~
make: *** [ffbuild/common.mak:67: libavcodec/libsvtav1.o] Error 1
make: *** Waiting for unfinished jobs....
```

**解決法**  
cross_compile_ffmpeg.shを書き換えてSVT-AV1のバージョンを下げる
```
  do_git_checkout https://gitlab.com/AOMediaCodec/SVT-AV1.git
　↓
  do_git_checkout https://gitlab.com/AOMediaCodec/SVT-AV1.git SVT-AV1_git v1.4.1
```
参考: http://www.neko.ne.jp/~freewing/software/windows_compile_ffmpeg_enable_fdk_aac/

## ビルドしたffmpegの入手
エクスプローラーに下記のアドレスを入力。
```
\\wsl$\Ubuntu
```
下記のパスにffmpegのexeがある。
```
home/<username>/ffmpeg-windows-build-helpers/sandbox/win64/ffmpeg_git_with_fdk_aac_n4.4.4/
```

あとは適当なフォルダに配置してPATHを通せばOK。