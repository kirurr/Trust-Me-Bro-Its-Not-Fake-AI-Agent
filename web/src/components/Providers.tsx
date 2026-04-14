import { ConfigProvider, theme } from "antd";
import { StyleProvider } from "@ant-design/cssinjs";
import { Provider } from "react-redux";
import { store } from "../app/store";

export default function Providers({ children }: { children: React.ReactNode }) {
  const configProps = { theme: {
		algorithm: theme.darkAlgorithm,
	} };
  return (
    <Provider store={store}>
      <StyleProvider layer>
        <ConfigProvider
					{...configProps}
				>{children}</ConfigProvider>
      </StyleProvider>
    </Provider>
  );
}
