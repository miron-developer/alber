const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function(app) {
    app.use(
        createProxyMiddleware("/api/", { target: "https://localhost:4330/", secure: false })
    );
    app.use(
        createProxyMiddleware("/sign/", { target: "https://localhost:4330/", secure: false })
    );
    app.use(
        createProxyMiddleware("/e/", { target: "https://localhost:4330/", secure: false })
    );
    app.use(
        createProxyMiddleware("/r/", { target: "https://localhost:4330/", secure: false })
    );
    app.use(
        createProxyMiddleware("/s/", { target: "https://localhost:4330/", secure: false })
    );
    app.use(
        createProxyMiddleware("/assets/", { target: "https://localhost:4330/", secure: false })
    );
};