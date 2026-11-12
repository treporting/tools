// The rules of eslint for every JavaScript package of every repo (SR_54). build.yml runs eslint with
// this file on every package.
//
// @implements SR_54@1.1.1
import js from "@eslint/js";
import globals from "globals";

export default [
  { ignores: ["**/dist/**", "**/node_modules/**", "**/coverage/**"] },
  js.configs.recommended,
  {
    files: ["**/*.{js,jsx}"],
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      globals: { ...globals.browser, ...globals.node },
      parserOptions: { ecmaFeatures: { jsx: true } },
    },
  },
];
