const test = require('node:test');
const assert = require('node:assert');
const nodemailer = require('nodemailer');

/**
 * Example tests using the Mailpit API.
 * This file demonstrates how to run email integration tests without external 
 * testing libraries (besides nodemailer already installed for sending).
 * 
 * We use:
 * - node:test (Node.js native test runner)
 * - fetch (Node.js 18+ native API)
 */

const MAILPIT_API = 'http://localhost:8025/api/v1';

// SMTP transport configuration for Mailpit
const transporter = nodemailer.createTransport({
  host: 'localhost',
  port: 1025,
  secure: false,
  auth: null
});

test('Mailpit API - Integration Test Demonstration', async (t) => {

  await t.test('1. Initial cleanup: should delete all messages', async () => {
    console.log('🧹 Cleaning inbox...');
    const response = await fetch(`${MAILPIT_API}/messages`, { method: 'DELETE' });
    assert.strictEqual(response.status, 200, 'API should return 200 on delete');
    
    const check = await fetch(`${MAILPIT_API}/messages`);
    const data = await check.json();
    assert.strictEqual(data.total, 0, 'Inbox should be empty');
  });

  await t.test('2. Sending: should send an email via SMTP', async () => {
    console.log('📧 Sending test email...');
    const info = await transporter.sendMail({
      from: 'api-test@example.com',
      to: 'you@example.com',
      subject: 'API Test Subject 🚀',
      text: 'This is an email sent via Node.js to be validated by the Mailpit API.',
      html: '<b>HTML content tested!</b>'
    });

    assert.ok(info.messageId, 'Should have a messageId after sending');
  });

  await t.test('3. Validation: should find the email via Mailpit API', async () => {
    // Wait a brief moment for Mailpit to process the message
    await new Promise(r => setTimeout(r, 300));

    console.log('🔍 Searching for email in API...');
    const response = await fetch(`${MAILPIT_API}/messages`);
    const data = await response.json();

    assert.strictEqual(data.total, 1, 'Should have exactly 1 message in the inbox');
    const msg = data.messages[0];

    assert.strictEqual(msg.Subject, 'API Test Subject 🚀');
    assert.strictEqual(msg.To[0].Address, 'you@example.com');
    assert.strictEqual(msg.From.Address, 'api-test@example.com');
    
    console.log(`✅ Email found! Mailpit ID: ${msg.ID}`);
  });

  await t.test('4. Details: should retrieve and validate HTML content', async () => {
    // First we get the list to obtain the current message ID
    const listRes = await fetch(`${MAILPIT_API}/messages`);
    const listData = await listRes.json();
    const messageId = listData.messages[0].ID;

    // We fetch the full message details
    const msgRes = await fetch(`${MAILPIT_API}/message/${messageId}`);
    const msgData = await msgRes.json();

    assert.ok(msgData.HTML.includes('HTML content tested!'), 'HTML content should contain the expected string');
    console.log('✅ HTML content validated successfully.');
    console.log('---');
  });

});
