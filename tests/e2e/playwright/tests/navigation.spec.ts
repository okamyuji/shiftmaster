import { test, expect, Page } from '@playwright/test';

// ログインヘルパー
async function login(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', email);
  await page.fill('input[type="password"], input[name="password"]', password);
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/');
}

test.describe('ナビゲーション', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
  });

  test('ダッシュボードが表示される', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
    await expect(page.locator('text=クイックアクション')).toBeVisible();
  });

  test('スタッフ一覧に遷移できる', async ({ page }) => {
    await page.click('a:has-text("スタッフ")');
    await expect(page).toHaveURL(/staffs/);
    await expect(page).toHaveTitle(/スタッフ/);
  });

  test('シフト種別一覧に遷移できる', async ({ page }) => {
    // ナビゲーションからシフトリンクをクリック
    await page.locator('nav a[href="/shifts"]').first().click();
    await expect(page).toHaveURL(/shifts/);
    await expect(page).toHaveTitle(/シフト/);
  });

  test('チーム一覧に遷移できる', async ({ page }) => {
    await page.click('a:has-text("チーム")');
    await expect(page).toHaveURL(/teams/);
    await expect(page).toHaveTitle(/チーム/);
  });

  test('勤務表一覧に遷移できる', async ({ page }) => {
    await page.click('a:has-text("勤務表")');
    await expect(page).toHaveURL(/schedules/);
    await expect(page).toHaveTitle(/勤務表/);
  });

  test('ユーザー管理に遷移できる（admin）', async ({ page }) => {
    await page.click('a:has-text("ユーザー管理")');
    await expect(page).toHaveURL(/admin\/users/);
    await expect(page).toHaveTitle(/ユーザー/);
  });
});

test.describe('super_admin専用機能', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'superadmin@example.com', 'SuperAdmin123$!');
  });

  test('組織セレクターが表示される', async ({ page }) => {
    // super_admin用の組織セレクターボタン
    const orgSelector = page.locator('button:has-text("組織を選択"), button:has-text("サンプル病院")');
    await expect(orgSelector).toBeVisible();
  });

  test('組織一覧が表示される', async ({ page }) => {
    // 組織セレクターをクリック
    await page.click('button:has-text("組織を選択"), button:has-text("サンプル病院")');

    // ドロップダウンに組織が表示される
    await expect(page.locator('a:has-text("サンプル病院")')).toBeVisible();
  });
});
