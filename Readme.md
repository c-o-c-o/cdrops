# CDrops
ごちゃまぜドロップスAPIを簡単に利用する為のアプリです  
aviutlへの読み込みをコマンドラインで実行できます  

# 使い方
```
cdrops.exe [DropParams]
DropParams
  *で区切られたパラメータを指定します
  レイヤー番号*シークの位置変更(ミリ秒)*ドロップするファイル*ドロップするファイル*ドロップする...
  ドロップするファイル は必ず絶対パスで指定してください
  パスは C:\\\\Users\\\\ユーザー名\\\\音声ファイル.wav の様な形式で指定するとより適切に処理されます
  ドロップするファイル は省略できます、その場合レイヤー番号に意味はなく、シーク位置の変更のみ行われます
  ドロップするファイル が存在する場合、シーク位置の変更はドロップしたファイルの最初か最後に固定されます(ごちゃまぜドロップスの仕様？)
  ドロップ後にシーク位置を動かしたい場合は、パラメータを分けてください
```

```
example
  - レイヤー1にドロップ後、オブジェクトの後ろから300ミリ秒シーク位置を動かす例
  cdrops.exe 1*1000*音声ファイル.wav*テキストファイル.txt 1*300

  - レイヤー3にドロップ後、オブジェクトの先から300ミリ秒シーク位置を動かす例
  cdrops.exe 3*0*音声ファイル.wav*テキストファイル.txt 1*300
```

# Licence
This software is released under the MIT License, see LICENSE.  

一部のファイルは "かんしくん" のファイルを改変して作成しています  
  かんしくんリポジトリ  
    https://github.com/oov/forcepser  
  改変元ファイル  
    https://github.com/oov/forcepser/blob/master/src/go/gcmz.go  
  改変後ファイル
    https://github.com/c-o-c-o/cdrops/blob/master/gcmz/gcmz.go