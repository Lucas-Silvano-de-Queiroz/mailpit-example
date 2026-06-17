const nodemailer = require('nodemailer');

const transporter = nodemailer.createTransport({
    host: 'localhost',
    port: 1025,
    secure: false,
    auth: null
});

async function send() {
    const info = await transporter.sendMail({
        from: 'sistema-node@exemplo.com',
        to: 'usuario-node@exemplo.com',
        subject: 'Teste Mailpit e Node.js 🚀',
        text: 'Olá! E-mail simplificado enviado via Node.js.',
    });

    console.log('E-mail enviado!', info.messageId);
}

send().catch(console.error);
