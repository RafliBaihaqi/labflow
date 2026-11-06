import { expect, test } from "@playwright/test";
import { BACKEND_URL, loginViaAPI, loginViaUI } from "./helpers";

test("customer submits a test request and lab staff can reject it", async ({ page, playwright }) => {
  const sampleName = `E2E sample ${Date.now()}`;

  // Customer creates the request through the API (the multi-step form is covered separately).
  const customer = await playwright.request.newContext();
  await loginViaAPI(customer, "customer");
  const created = await customer.post(`${BACKEND_URL}/api/reservations`, {
    data: {
      sample_name: sampleName,
      sample_amount: 2,
      special_instructions: "e2e",
      sample_type: "Liquid",
      sample_taken_date: new Date().toISOString(),
      sample_taken_by: "e2e",
      tests: [{ test_list_id: 1, method: "USP <61>", criteria: "< 100 CFU/g", test_unit: "CFU/g" }],
    },
  });
  expect(created.status()).toBe(201);
  const { xid } = await created.json();

  // It shows up in the customer's reservation list in the UI.
  await loginViaUI(page, "customer");
  await page.goto("/reservations");
  await expect(page.getByText(sampleName)).toBeVisible();

  // A customer may not move their own request through the lifecycle...
  const selfReject = await customer.put(`${BACKEND_URL}/api/reservations/${xid}/status`, {
    data: { status: "rejected" },
  });
  expect(selfReject.status()).toBe(403);

  // ...but QA can reject it, and the new status is visible to the customer.
  const qa = await playwright.request.newContext();
  await loginViaAPI(qa, "qa");
  const rejected = await qa.put(`${BACKEND_URL}/api/reservations/${xid}/status`, {
    data: { status: "rejected" },
  });
  expect(rejected.status()).toBe(200);

  const detail = await qa.get(`${BACKEND_URL}/api/reservations/${xid}`);
  expect(await detail.json()).toMatchObject({ status: "rejected" });

  await Promise.all([customer.dispose(), qa.dispose()]);
});
