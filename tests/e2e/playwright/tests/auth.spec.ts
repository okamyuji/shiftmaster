import { test, expect } from '@playwright/test';

test.describe('認証機能', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
  });

  test('ログインページが正常に表示される', async ({ page }) => {
    await expect(page).toHaveTitle(/ログイン/);
    await expect(page.locator('h1, h2').filter({ hasText: 'ログイン' })).toBeVisible();
    await expect(page.locator('input[type="email"], input[name="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"], input[name="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('正常なログイン - admin', async ({ page }) => {
    await page.fill('input[type="email"], input[name="email"]', 'admin@example.com');
    await page.fill('input[type="password"], input[name="password"]', 'Password123$!');
    await page.click('button[type="submit"]');

    // ダッシュボードにリダイレクト
    await expect(page).toHaveURL('/');
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('正常なログイン - super_admin', async ({ page }) => {
    await page.fill('input[type="email"], input[name="email"]', 'superadmin@example.com');
    await page.fill('input[type="password"], input[name="password"]', 'SuperAdmin123$!');
    await page.click('button[type="submit"]');

    // ダッシュボードにリダイレクト
    await expect(page).toHaveURL('/');
    await expect(page).toHaveTitle(/ダッシュボード/);

    // super_admin用の組織セレクターが表示される
    await expect(page.getByRole('button', { name: '組織を選択' }).or(page.getByRole('button', { name: /サンプル病院/ }))).toBeVisible();
  });

  test('不正なパスワードでログイン失敗', async ({ page }) => {
    await page.fill('input[type="email"], input[name="email"]', 'admin@example.com');
    await page.fill('input[type="password"], input[name="password"]', 'WrongPassword!');
    await page.click('button[type="submit"]');

    // ログインページに留まる
    await expect(page).toHaveURL(/login/);
  });

  test('存在しないユーザーでログイン失敗', async ({ page }) => {
    await page.fill('input[type="email"], input[name="email"]', 'notexist@example.com');
    await page.fill('input[type="password"], input[name="password"]', 'Password123$!');
    await page.click('button[type="submit"]');

    // ログインページに留まる
    await expect(page).toHaveURL(/login/);
  });

  test('ログアウト', async ({ page }) => {
    // ログイン
    await page.fill('input[type="email"], input[name="email"]', 'admin@example.com');
    await page.fill('input[type="password"], input[name="password"]', 'Password123$!');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('/');

    // ユーザーメニューを開いてログアウト
    await page.click('button:has(svg):has-text("A"), button:has-text("A")');
    await page.click('button:has-text("ログアウト")');

    // ログインページにリダイレクト
    await expect(page).toHaveURL(/login/);
  });
});

test.describe('認証が必要なページへのアクセス', () => {
  test.use({ storageState: { cookies: [], origins: [] } });

  test('未認証でダッシュボードアクセス時は認証エラーまたはログインにリダイレクト', async ({ page }) => {
    const response = await page.goto('/');
    // 401 Unauthorizedまたはログインへのリダイレクト
    const status = response?.status();
    const url = page.url();
    expect(status === 401 || url.includes('login')).toBeTruthy();
  });

  test('未認証でスタッフ一覧アクセス時は認証エラーまたはログインにリダイレクト', async ({ page }) => {
    const response = await page.goto('/staffs');
    const status = response?.status();
    const url = page.url();
    expect(status === 401 || url.includes('login')).toBeTruthy();
  });

  test('未認証で管理者ページアクセス時は認証エラーまたはログインにリダイレクト', async ({ page }) => {
    const response = await page.goto('/admin/users');
    const status = response?.status();
    const url = page.url();
    expect(status === 401 || url.includes('login')).toBeTruthy();
  });
});
