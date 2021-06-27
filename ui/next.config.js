const productionEnv = {
  serverPrefix: "",
};

const developmentEnv = {
  serverPrefix: "http://localhost:8080",
};

const env =
  process.env.NODE_ENV === "production" ? productionEnv : developmentEnv;

module.exports = {
  reactStrictMode: true,
  env,
};
