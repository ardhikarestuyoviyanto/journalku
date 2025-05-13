import prettierPlugin from "eslint-plugin-prettier";
import eslintPluginReact from "eslint-plugin-react";
import eslintPluginReactHooks from "eslint-plugin-react-hooks";
import globals from "globals";

export default [
  { ignores: ["eslint.config.mjs"] },
  {
    ...eslintPluginReact.configs.flat.recommended,
    ...eslintPluginReact.configs.flat["jsx-runtime"],
    files: ["src/**/*.{js,jsx}"],
    plugins: {
      eslintPluginReact,
      "react-hooks": eslintPluginReactHooks,
    },
    rules: {
      "prettier/prettier": ["error", { endOfLine: "lf" }],
    },
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
      ecmaVersion: "latest",
      sourceType: "module",
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
    },
    settings: {
      react: {
        version: "detect",
      },
    },
    rules: {
      ...eslintPluginReactHooks.configs.recommended.rules,
    },
  },
  eslintPluginPrettierRecommended,
];
