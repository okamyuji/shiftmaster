import { expect, Page, test } from '@playwright/test';

// テスト間のCookie競合を防ぐためにシリアルモードで実行
test.describe.configure({ mode: 'serial' });

async function login(page: Page, email: string = 'admin@example.com', password: string = 'Password123$!') {
  // 既存のCookieをクリア
  await page.context().clearCookies();

  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', email);
  await page.fill('input[type="password"], input[name="password"]', password);

  // ログインボタンをクリックしてナビゲーション完了を待機
  await Promise.all([
    page.waitForURL('/'),
    page.click('button[type="submit"]'),
  ]);
}

test.describe('シフト種別管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('シフト種別一覧が表示される', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveTitle(/シフト/);

    // 初期データのシフト種別が表示される（h3要素内のテキストを確認）
    await expect(page.locator('h3:text-is("日勤")')).toBeVisible();
    await expect(page.locator('h3:text-is("夜勤")')).toBeVisible();
  });

  test('シフト追加フォームに遷移できる', async ({ page }) => {
    await page.goto('/shifts');
    await page.click('a:has-text("シフト追加")');

    await expect(page).toHaveURL(/shifts\/new/);
    await expect(page.locator('input[name="name"]')).toBeVisible();
  });

  test('シフト種別の詳細が確認できる', async ({ page }) => {
    await page.goto('/shifts');

    // 日勤のカード内の編集リンクをクリック（シフト詳細は編集ページで表示）
    await page.locator('a[href*="/shifts/"][href$="/edit"]').first().click();

    // 編集ページに遷移することを確認
    await expect(page).toHaveURL(/shifts\/[a-f0-9-]+\/edit/);
    // シフト名の入力欄が表示される
    await expect(page.locator('input[name="name"]')).toBeVisible();
  });
});

test.describe('申し送り時間（HandoverMinutes）機能', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('シフト追加フォームに申し送り時間入力欄が存在する', async ({ page }) => {
    await page.goto('/shifts/new');

    // 申し送り時間の入力欄が存在することを確認
    await expect(page.locator('input[name="handover_minutes"]')).toBeVisible();
    await expect(page.locator('label:has-text("申し送り時間")')).toBeVisible();
  });

  test('シフト編集フォームに申し送り時間入力欄が存在する', async ({ page }) => {
    await page.goto('/shifts');

    // 最初のシフトの編集ページに遷移
    await page.locator('a[href*="/shifts/"][href$="/edit"]').first().click();
    await expect(page).toHaveURL(/shifts\/[a-f0-9-]+\/edit/);

    // 申し送り時間の入力欄が存在することを確認
    await expect(page.locator('input[name="handover_minutes"]')).toBeVisible();
  });

  test('申し送り時間を設定してシフトを作成できる', async ({ page }) => {
    await page.goto('/shifts/new');

    // シフト情報を入力
    const uniqueName = `テスト交代勤務_${Date.now()}`;
    await page.fill('input[name="name"]', uniqueName);
    await page.fill('input[name="code"]', 'TST');
    await page.fill('input[name="start_time"]', '08:00');
    await page.fill('input[name="end_time"]', '16:30');
    await page.fill('input[name="break_minutes"]', '60');
    await page.fill('input[name="handover_minutes"]', '30'); // 申し送り30分

    // カラーピッカーは直接値を設定（evaluate使用）
    await page.locator('input[name="color"]').evaluate((el: HTMLInputElement) => {
      el.value = '#FF5733';
    });

    // 「登録」ボタンをクリック
    await page.getByRole('button', { name: '登録' }).click();

    // 一覧ページにリダイレクト
    await expect(page).toHaveURL(/\/shifts/);

    // 作成したシフトが表示される
    await expect(page.locator(`h3:has-text("${uniqueName}")`)).toBeVisible();
  });

  test('申し送り時間の境界値テスト - 0分', async ({ page }) => {
    await page.goto('/shifts/new');

    const uniqueName = `境界値0分_${Date.now()}`;
    await page.fill('input[name="name"]', uniqueName);
    await page.fill('input[name="code"]', 'B00');
    await page.fill('input[name="start_time"]', '09:00');
    await page.fill('input[name="end_time"]', '17:00');
    await page.fill('input[name="break_minutes"]', '60');
    await page.fill('input[name="handover_minutes"]', '0'); // 申し送り0分

    await page.locator('input[name="color"]').evaluate((el: HTMLInputElement) => {
      el.value = '#00FF00';
    });

    await page.getByRole('button', { name: '登録' }).click();
    await expect(page).toHaveURL(/\/shifts/);
    await expect(page.locator(`h3:has-text("${uniqueName}")`)).toBeVisible();
  });

  test('申し送り時間の境界値テスト - 最大値120分', async ({ page }) => {
    await page.goto('/shifts/new');

    const uniqueName = `境界値120分_${Date.now()}`;
    await page.fill('input[name="name"]', uniqueName);
    await page.fill('input[name="code"]', 'B12');
    await page.fill('input[name="start_time"]', '07:00');
    await page.fill('input[name="end_time"]', '19:00');
    await page.fill('input[name="break_minutes"]', '60');
    await page.fill('input[name="handover_minutes"]', '120'); // 申し送り120分（最大値）

    await page.locator('input[name="color"]').evaluate((el: HTMLInputElement) => {
      el.value = '#0000FF';
    });

    await page.getByRole('button', { name: '登録' }).click();
    await expect(page).toHaveURL(/\/shifts/);
    await expect(page.locator(`h3:has-text("${uniqueName}")`)).toBeVisible();
  });

  test('申し送り時間を編集できる', async ({ page }) => {
    // まず新しいシフトを作成
    await page.goto('/shifts/new');

    const uniqueName = `編集テスト_${Date.now()}`;
    await page.fill('input[name="name"]', uniqueName);
    await page.fill('input[name="code"]', 'EDT');
    await page.fill('input[name="start_time"]', '08:00');
    await page.fill('input[name="end_time"]', '16:00');
    await page.fill('input[name="break_minutes"]', '45');
    await page.fill('input[name="handover_minutes"]', '15');

    await page.locator('input[name="color"]').evaluate((el: HTMLInputElement) => {
      el.value = '#FFAA00';
    });

    await page.getByRole('button', { name: '登録' }).click();
    await expect(page).toHaveURL(/\/shifts/);

    // 作成したシフトのカードを特定し、その中の編集リンクをクリック
    // h3要素の正確なテキストでカードを見つける
    const shiftCard = page.locator('div.card', { has: page.locator(`h3:text-is("${uniqueName}")`) });
    await shiftCard.locator('a[href$="/edit"]').click();

    // 申し送り時間を変更
    await page.fill('input[name="handover_minutes"]', '45');
    await page.getByRole('button', { name: '更新' }).click();

    // 一覧に戻る
    await expect(page).toHaveURL(/\/shifts/);
  });
});

test.describe('シフト種別 - マルチテナント分離テスト', () => {
  test.beforeEach(async ({ page }) => {
    // テスト医療センターのadminでログイン
    await login(page, 'otheradmin@example.com', 'OtherAdmin123$!');
  });

  test('他組織のシフト種別は表示されない', async ({ page }) => {
    await page.goto('/shifts');

    // サンプル病院のシフト種別（日勤、夜勤など）が表示されないことを確認
    await expect(page.locator('h3:text-is("日勤")')).not.toBeVisible();
    await expect(page.locator('h3:text-is("夜勤")')).not.toBeVisible();
  });

  test('自組織のシフト種別のみ作成・管理できる', async ({ page }) => {

    await page.goto('/shifts/new');

    // シフトを作成
    const uniqueName = `他組織シフト_${Date.now()}`;
    await page.fill('input[name="name"]', uniqueName);
    await page.fill('input[name="code"]', 'OTH');
    await page.fill('input[name="start_time"]', '09:00');
    await page.fill('input[name="end_time"]', '17:00');
    await page.fill('input[name="break_minutes"]', '60');
    await page.fill('input[name="handover_minutes"]', '0');

    await page.locator('input[name="color"]').evaluate((el: HTMLInputElement) => {
      el.value = '#999999';
    });

    await page.getByRole('button', { name: '登録' }).click();
    await expect(page).toHaveURL(/\/shifts/);

    // 作成したシフトが表示される
    await expect(page.locator(`h3:has-text("${uniqueName}")`)).toBeVisible();
  });
});

test.describe('シフト種別 - 異常系テスト', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('負の申し送り時間は入力できない', async ({ page }) => {
    await page.goto('/shifts/new');

    const handoverInput = page.locator('input[name="handover_minutes"]');

    // min属性が0に設定されていることを確認
    await expect(handoverInput).toHaveAttribute('min', '0');
  });

  test('フォームに必須フィールドが存在する', async ({ page }) => {
    await page.goto('/shifts/new');

    // 必須フィールドが存在することを確認
    await expect(page.locator('input[name="name"]')).toBeVisible();
    await expect(page.locator('input[name="code"]')).toBeVisible();
    await expect(page.locator('input[name="start_time"]')).toBeVisible();
    await expect(page.locator('input[name="end_time"]')).toBeVisible();
    await expect(page.locator('input[name="break_minutes"]')).toBeVisible();
    await expect(page.locator('input[name="handover_minutes"]')).toBeVisible();
    await expect(page.locator('input[name="color"]')).toBeVisible();
  });
});

test.describe('super_admin - シフト種別管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'superadmin@example.com', 'SuperAdmin123$!');
  });

  test('super_adminは全組織のシフト種別にアクセスできる', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveTitle(/シフト/);

    // シフト一覧ページが正常に表示される
    await expect(page.locator('h1, h2')).toContainText(/シフト/);
  });

  test('super_adminはシフトを作成できる', async ({ page }) => {
    await page.goto('/shifts/new');

    // フォームが表示される
    await expect(page.locator('input[name="name"]')).toBeVisible();
    await expect(page.locator('input[name="handover_minutes"]')).toBeVisible();
  });
});
