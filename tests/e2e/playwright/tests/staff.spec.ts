import { test, expect, Page } from '@playwright/test';

async function login(page: Page) {
  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', 'admin@example.com');
  await page.fill('input[type="password"], input[name="password"]', 'Password123$!');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/');
}

test.describe('スタッフ管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('スタッフ一覧が表示される', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveTitle(/スタッフ/);

    // 初期データの山田花子が表示される（リンク要素を確認）
    await expect(page.locator('a:has-text("山田 花子")')).toBeVisible();
  });

  test('スタッフ検索ができる', async ({ page }) => {
    await page.goto('/staffs');

    // 検索ボックスに入力
    await page.fill('input[placeholder*="検索"]', '山田');

    // 検索結果が表示される（リンク要素を確認）
    await expect(page.locator('a:has-text("山田 花子")')).toBeVisible();
  });

  test('スタッフ追加フォームに遷移できる', async ({ page }) => {
    await page.goto('/staffs');
    await page.click('a:has-text("スタッフ追加")');

    await expect(page).toHaveURL(/staffs\/new/);
    await expect(page.locator('input[name="first_name"], input[name="firstName"]')).toBeVisible();
  });

  test('スタッフ詳細が表示される', async ({ page }) => {
    await page.goto('/staffs');

    // スタッフ名をクリック
    await page.locator('a:has-text("山田 花子")').first().click();

    // 詳細ページに遷移することを確認
    await expect(page).toHaveURL(/staffs\/[a-f0-9-]+/);
    await expect(page.locator('h1')).toContainText('山田');
  });
});
