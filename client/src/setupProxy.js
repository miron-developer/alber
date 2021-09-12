const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function(app) {
    app.use(
        createProxyMiddleware("/api/", { target: "http://localhost:4330/" })
    );
    app.use(
        createProxyMiddleware("/sign/", { target: "http://localhost:4330/" })
    );
    app.use(
        createProxyMiddleware("/e/", { target: "http://localhost:4330/" })
    );
    app.use(
        createProxyMiddleware("/s/", { target: "http://localhost:4330/" })
    );
    app.use(
        createProxyMiddleware("/assets/", { target: "http://localhost:4330/" })
    );
};