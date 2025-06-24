# Claude Codeの使い方

## 初期設定

1. Terminal（`Ctrl+Shift+@`）を開き、`claude --dangerously-skip-permissions`を実行します  
   claudeの初期設定を行ってください  
   `--dangerously-skip-permissions` はclaudeに全ての操作を許可します（Dev Container上で実行するため許容します）
  
2. その他の使い方は以下を参照してください
   https://docs.anthropic.com/ja/docs/claude-code/cli-usage

## 本リポジトリのClaude Code活用方針

* DBのスキーマ定義（ent/schema）とOpenAPI定義（swagger.yml）を自身の手で作成し、その実装をClaude Codeに任せる
  * 上記のファイル修正後、`現在の変更状態を確認して、影響する部分を修正して` と指示するとAPIを実装してくれる
  * `POST /users/{id}/todosを実装して` でも可
 
## Claude Codeの便利機能

