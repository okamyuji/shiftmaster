import { expect, Page, test } from '@playwright/test';

// ログインヘルパー
async function login(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.fill('input[type="email"], input[name="email"]', email);
  await page.fill('input[type="password"], input[name="password"]', password);
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/');
}

test.describe('認可テスト - super_admin', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'superadmin@example.com', 'SuperAdmin123$!');
  });

  test('ダッシュボードにアクセスできる', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('ユーザー管理にアクセスできる', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveURL(/admin\/users/);
    await expect(page).toHaveTitle(/ユーザー/);
  });

  test('スタッフ一覧にアクセスできる', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
  });

  test('シフト種別一覧にアクセスできる', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveURL(/shifts/);
  });

  test('チーム一覧にアクセスできる', async ({ page }) => {
    await page.goto('/teams');
    await expect(page).toHaveURL(/teams/);
  });

  test('勤務表一覧にアクセスできる', async ({ page }) => {
    await page.goto('/schedules');
    await expect(page).toHaveURL(/schedules/);
  });

  test('組織セレクターが表示される', async ({ page }) => {
    await expect(page.locator('button:has-text("組織を選択")')).toBeVisible();
  });

  test('ロール表示が全体管理者', async ({ page }) => {
    // ユーザーメニューを開く
    await page.click('button:has-text("U"), button:has-text("S")');
    await expect(page.locator('p:has-text("全体管理者")')).toBeVisible();
  });
});

test.describe('認可テスト - admin（サンプル病院）', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'admin@example.com', 'Password123$!');
  });

  test('ダッシュボードにアクセスできる', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('ユーザー管理にアクセスできる', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveURL(/admin\/users/);
    await expect(page).toHaveTitle(/ユーザー/);
  });

  test('スタッフ一覧にアクセスできる', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
  });

  test('組織名が表示される', async ({ page }) => {
    // ヘッダーまたはサイドバーに組織名が表示される
    await expect(page.locator('span.text-xs:has-text("サンプル病院")')).toBeVisible();
  });

  test('ロール表示がテナント管理者', async ({ page }) => {
    // ユーザーメニューを開く
    await page.click('button:has-text("U"), button:has-text("A")');
    await expect(page.locator('p:has-text("テナント管理者")')).toBeVisible();
  });

  test('組織セレクターが表示されない（super_admin専用）', async ({ page }) => {
    await expect(page.locator('button:has-text("組織を選択")')).not.toBeVisible();
  });

  test('自組織のスタッフデータが見える（マルチテナント分離）', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
    // サンプル病院のスタッフ（山田花子）が見えることを確認
    await expect(page.getByRole('link', { name: '山田 花子' })).toBeVisible();
  });

  test('自組織のシフト種別データが見える（マルチテナント分離）', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveURL(/shifts/);
    // サンプル病院のシフト種別（日勤）が見えることを確認
    await expect(page.getByRole('heading', { name: '日勤', exact: true })).toBeVisible();
  });

  test('自組織のチームデータが見える（マルチテナント分離）', async ({ page }) => {
    await page.goto('/teams');
    await expect(page).toHaveURL(/teams/);
    // サンプル病院のチーム（看護1チーム）が見えることを確認
    await expect(page.locator('text=看護1チーム')).toBeVisible();
  });

  test('自組織のユーザーが見える（マルチテナント分離）', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveURL(/admin\/users/);
    // サンプル病院のユーザー（admin@example.com）が見えることを確認
    await expect(page.locator('text=admin@example.com')).toBeVisible();
  });
});

test.describe('認可テスト - manager', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'manager@example.com', 'Manager123$!');
  });

  test('ダッシュボードにアクセスできる', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('ユーザー管理にアクセスすると403エラー', async ({ page }) => {
    const response = await page.goto('/admin/users');
    expect(response?.status()).toBe(403);
  });

  test('スタッフ一覧にアクセスできる', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
  });

  test('シフト種別一覧にアクセスできる', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveURL(/shifts/);
  });

  test('勤務表一覧にアクセスできる', async ({ page }) => {
    await page.goto('/schedules');
    await expect(page).toHaveURL(/schedules/);
  });

  test('ロール表示がマネージャー', async ({ page }) => {
    // ユーザーメニューを開く
    await page.click('button:has-text("U"), button:has-text("M")');
    await expect(page.locator('p:has-text("マネージャー")')).toBeVisible();
  });
});

test.describe('認可テスト - user（一般ユーザー）', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'user@example.com', 'User123$!');
  });

  test('ダッシュボードにアクセスできる', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('ユーザー管理にアクセスすると403エラー', async ({ page }) => {
    const response = await page.goto('/admin/users');
    expect(response?.status()).toBe(403);
  });

  test('スタッフ一覧にアクセスできる', async ({ page }) => {
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
  });

  test('シフト種別一覧にアクセスできる', async ({ page }) => {
    await page.goto('/shifts');
    await expect(page).toHaveURL(/shifts/);
  });

  test('ロール表示が一般ユーザー', async ({ page }) => {
    // ユーザーメニューを開く
    await page.click('button:has-text("U"), button:has-text("H")');
    await expect(page.locator('p:has-text("一般ユーザー")')).toBeVisible();
  });
});

test.describe('認可テスト - 他組織admin（テスト医療センター）', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, 'otheradmin@example.com', 'OtherAdmin123$!');
  });

  test('ダッシュボードにアクセスできる', async ({ page }) => {
    await expect(page).toHaveTitle(/ダッシュボード/);
  });

  test('ユーザー管理にアクセスできる', async ({ page }) => {
    await page.goto('/admin/users');
    await expect(page).toHaveURL(/admin\/users/);
  });

  test('組織名が表示される', async ({ page }) => {
    // ヘッダーに組織名が表示される
    await expect(page.locator('span.text-xs:has-text("テスト医療センター")')).toBeVisible();
  });

  test('サンプル病院のスタッフデータは見えない（マルチテナント分離）', async ({ page }) => {
    // スタッフ一覧を表示
    await page.goto('/staffs');
    await expect(page).toHaveURL(/staffs/);
    // サンプル病院のスタッフ（山田花子）が見えないことを確認
    // テスト医療センターにはスタッフが登録されていないため、リンクとして存在しないはず
    const yamadaCount = await page.locator('a:has-text("山田")').count();
    expect(yamadaCount).toBe(0);
  });

  test('サンプル病院のシフト種別データは見えない（マルチテナント分離）', async ({ page }) => {
    // シフト種別一覧を表示
    await page.goto('/shifts');
    await expect(page).toHaveURL(/shifts/);
    // サンプル病院のシフト種別（日勤など）が見えないことを確認
    // テスト医療センターにはシフト種別が登録されていない場合、存在しないはず
    const shiftCards = await page.locator('div.card h3:has-text("日勤")').count();
    expect(shiftCards).toBe(0);
  });

  test('サンプル病院のチームデータは見えない（マルチテナント分離）', async ({ page }) => {
    // チーム一覧を表示
    await page.goto('/teams');
    await expect(page).toHaveURL(/teams/);
    // サンプル病院のチーム（看護1チームなど）が見えないことを確認
    const teamCount = await page.locator('a:has-text("看護1チーム")').count();
    expect(teamCount).toBe(0);
  });

  test('サンプル病院のユーザーは見えない（マルチテナント分離）', async ({ page }) => {
    // ユーザー一覧を表示
    await page.goto('/admin/users');
    await expect(page).toHaveURL(/admin\/users/);
    // サンプル病院のユーザー（admin@example.com）が見えないことを確認
    // 注: otheradmin@example.comには'admin@example.com'が含まれるため、完全一致で検索
    const adminCount = await page.locator('text="admin@example.com"').count();
    expect(adminCount).toBe(0);
    // テスト医療センターのユーザー（otheradmin@example.com）は見える
    await expect(page.locator('text=otheradmin@example.com')).toBeVisible();
  });
});

test.describe('境界値・エッジケーステスト', () => {
  test('未認証でダッシュボードにアクセスすると401', async ({ page }) => {
    const response = await page.goto('/');
    expect(response?.status()).toBe(401);
  });

  test('未認証でスタッフ一覧にアクセスすると401', async ({ page }) => {
    const response = await page.goto('/staffs');
    expect(response?.status()).toBe(401);
  });

  test('未認証でユーザー管理にアクセスすると401', async ({ page }) => {
    const response = await page.goto('/admin/users');
    expect(response?.status()).toBe(401);
  });

  test('ログアウト後に保護されたページにアクセスできない', async ({ page }) => {
    // ログイン
    await login(page, 'admin@example.com', 'Password123$!');
    await expect(page).toHaveTitle(/ダッシュボード/);

    // ログアウト
    await page.click('button:has-text("A"), button:has-text("U")');
    await page.click('button:has-text("ログアウト")');
    await expect(page).toHaveURL(/login/);

    // 保護されたページにアクセス
    const response = await page.goto('/staffs');
    expect(response?.status()).toBe(401);
  });
});
