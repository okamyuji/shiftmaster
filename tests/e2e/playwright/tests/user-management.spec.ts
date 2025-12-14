import { test, expect, Page } from '@playwright/test';

async function loginAsAdmin(page: Page) {
  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', 'admin@example.com');
  await page.fill('input[type="password"], input[name="password"]', 'Password123$!');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/');
}

async function loginAsSuperAdmin(page: Page) {
  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', 'superadmin@example.com');
  await page.fill('input[type="password"], input[name="password"]', 'SuperAdmin123$!');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/');
}

test.describe('ユーザー管理 - admin', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test('ユーザー一覧が表示される', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveTitle(/ユーザー/);

    // 自分のユーザーが表示される（p要素内のメールを確認）
    await expect(page.locator('p:text-is("admin@example.com")')).toBeVisible();
  });

  test('ユーザー追加フォームに遷移できる', async ({ page }) => {
    await page.goto('/admin/users');
    await page.click('a:has-text("ユーザー追加")');

    await expect(page).toHaveURL(/admin\/users\/new/);
    await expect(page.locator('input[name="email"]')).toBeVisible();
  });

  test('ロール選択肢にsuper_adminが含まれない', async ({ page }) => {
    await page.goto('/admin/users/new');

    // ロール選択肢を開く
    const roleSelect = page.locator('select[name="role"]');

    // super_adminが選択肢にない
    const options = await roleSelect.locator('option').allTextContents();
    expect(options).not.toContain('全体管理者');
  });
});

test.describe('ユーザー管理 - super_admin', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsSuperAdmin(page);
  });

  test('ユーザー一覧が表示される', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveTitle(/ユーザー/);
  });

  test('ロール選択肢にsuper_adminが含まれる', async ({ page }) => {
    await page.goto('/admin/users/new');

    // ロール選択肢を開く
    const roleSelect = page.locator('select[name="role"]');

    // super_adminが選択肢にある（テキストをトリムして比較）
    const options = await roleSelect.locator('option').allTextContents();
    const trimmedOptions = options.map(opt => opt.trim());
    expect(trimmedOptions).toContain('全体管理者');
  });
});
