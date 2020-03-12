# miniSchemeGo

miniSchemeGoとは、Goで実装したとても小さなSchemeインタプリタです。
技術書典8頒布物「虎の穴ラボ TECH BOOK vol.1」の「Go 言語で作る処理系実装入門」記事のためのサンプルコードとなります。

## 仕様

以下機能が使えます。

 - car
 - cdr
 - cons
 - lambda
 - if
 - define
 - 比較演算(<,>,<=,>=)
 - 足し算、引き算
 - quote関数の簡略文字'(シングルクォート)

## ビルドと実行

```
$ go build
$ ./miniSchemeGo
miniSchemeGo>
```

`Control + C` で終了できます。

## サンプル

 - 階乗計算

入力（階乗を出力する関数定義）
```
(define kai (lambda (x) (if (> x 1) (* x (kai (- x 1))) 1)))
```

入力：10の階乗を計算
```
(kai 10)
```
出力：3628800

 - リスト処理

```
(car (cdr (cdr '(a b c d e))))
```

出力：c
