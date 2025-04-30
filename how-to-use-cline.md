# Clineの使い方

## 初期設定

1. コマンドパレット（`Ctrl+Shift+P`）を開き、`Cline: Open In New Tab`を実行します
  ![image](https://github.com/user-attachments/assets/c2d617ad-2c55-4f10-b985-31f0b62d2f23)  
  
2. 設定を開いてAPIキーの入力とモデルを選択します
   ![image](https://github.com/user-attachments/assets/fbdd8f6c-3999-4652-8c59-52d1a26ad0b7)

3. プラスをクリックして新規チャットを開きます
   ![image](https://github.com/user-attachments/assets/6053841b-9c27-4f90-a8cf-19c9122dcb9e)

4. Auto-approveをクリックします
   ![image](https://github.com/user-attachments/assets/9f16e6cf-aa3f-4aad-89d7-a1060d15f1df)

5. 全てにチェックを入れます（全ての操作を確認なしで実行しますが、devcontainer上なのでコンテナ外には影響しません）
   ![image](https://github.com/user-attachments/assets/cebabd14-f48d-4b3b-8293-b0aa9e40f907)

6. Auto-approveをクリックして閉じます
   ![image](https://github.com/user-attachments/assets/d45b1409-8063-432c-9f27-d37fdec26867)

7. 赤枠でClineに指示を出せます
   サンプルとして「POST /users/{id}/todosを実装して」と伝えると実装してくれるはずです  
   ![image](https://github.com/user-attachments/assets/70607b75-6e44-49e8-8a2c-01ea66829d2c)

## 本リポジトリのCline活用方針

* DBのスキーマ定義（ent/schema）とOpenAPI定義（swagger.yml）を自身の手で作成し、その実装をClineに任せる
  * 上記のファイル修正後、`現在の変更状態を確認して、影響する部分を修正して` と指示するとAPIを実装してくれる
  * `POST /users/{id}/todosを実装して` でも可
 
## Clineの便利機能

