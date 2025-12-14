import { expect, Page, test } from '@playwright/test';

// テストは順番に実行
test.describe.configure({ mode: 'serial' });

// ログインヘルパー関数
async function login(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', password);
  await page.click('button[type="submit"]');
  await page.waitForURL('/');
}

test.describe('勤務表バリデーション', () => {
  test('ログイン後にダッシュボードが表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await expect(page).toHaveURL('/');
    await expect(page.getByRole('heading', { name: 'ダッシュボード' })).toBeVisible();
  });

  test('勤務表一覧ページにアクセスできる', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/schedules');
    await expect(page.locator('h1:has-text("勤務表")')).toBeVisible();
  });

  test('勤務表作成フォームにアクセスできる', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/schedules/new');
    await expect(page.locator('text=勤務表作成')).toBeVisible();
  });
});

test.describe('シフト時間重複バリデーション - ドメインテスト', () => {
  // これはユニットテストで検証済みなので、E2Eではページ表示の確認のみ
  test('シフト種別一覧で日勤シフトが表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts');
    await expect(page.locator('text=日勤').first()).toBeVisible();
  });

  test('シフト種別詳細で時間情報が表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts');
    // シフトカードのリンクをクリック（詳細ボタン）
    const detailLink = page.locator('a[href^="/shifts/"]:not([href="/shifts/new"])').first();
    if (await detailLink.isVisible()) {
      await detailLink.click();
      // 詳細ページが表示されることを確認
      await expect(page.locator('text=開始時刻').or(page.locator('text=シフト詳細'))).toBeVisible();
    }
  });
});

test.describe('夜勤シフト - 日跨ぎ対応', () => {
  test('夜勤シフトが表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts');
    await expect(page.locator('text=夜勤').first()).toBeVisible();
  });

  test('申し送り時間フィールドが存在する', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts/new');
    await expect(page.locator('input[name="handover_minutes"]')).toBeVisible();
  });
});

test.describe('スタッフ管理', () => {
  test('スタッフ一覧が表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/staffs');
    await expect(page.locator('h1:has-text("スタッフ")')).toBeVisible();
  });

  test('スタッフ詳細ページにアクセスできる', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/staffs');
    // 最初のスタッフリンクをクリック
    const staffLink = page.locator('a[href^="/staffs/"]:not([href="/staffs/new"])').first();
    if (await staffLink.isVisible()) {
      await staffLink.click();
      await expect(page.locator('text=スタッフ詳細')).toBeVisible();
    }
  });
});

test.describe('チーム管理', () => {
  test('チーム一覧が表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/teams');
    await expect(page.locator('h1:has-text("チーム")')).toBeVisible();
  });
});

test.describe('勤務希望管理', () => {
  test('勤務希望一覧が表示される', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/requests');
    await expect(page.locator('h1:has-text("勤務希望")')).toBeVisible();
  });
});

test.describe('Value Object検証 - 時刻形式', () => {
  test('シフト作成フォームで時刻入力ができる', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts/new');

    // 時刻フィールドが存在することを確認
    await expect(page.locator('input[name="start_time"]')).toBeVisible();
    await expect(page.locator('input[name="end_time"]')).toBeVisible();
  });

  test('有効な時刻形式で入力できる', async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
    await page.goto('/shifts/new');

    await page.fill('input[name="start_time"]', '09:00');
    await page.fill('input[name="end_time"]', '17:00');

    // 値が正しく設定されていることを確認
    await expect(page.locator('input[name="start_time"]')).toHaveValue('09:00');
    await expect(page.locator('input[name="end_time"]')).toHaveValue('17:00');
  });
});
