import { ConfigProvider, theme } from "antd";
import { Provider } from "react-redux";
import { store } from "../app/store";

export default function Providers({ children }: { children: React.ReactNode }) {
  const configProps = { theme: { algorithm: theme.darkAlgorithm } };
  return (
    <Provider store={store}>
      <ConfigProvider {...configProps}>{children}</ConfigProvider>
    </Provider>
  );
}
