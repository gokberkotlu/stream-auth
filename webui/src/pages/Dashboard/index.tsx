import React from "react";
import {
  LaptopOutlined,
  NotificationOutlined,
  UserOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import type { MenuProps } from "antd";
import { Breadcrumb, Layout, Menu, theme } from "antd";
import "./styles.css";
import { useNavigate } from "react-router-dom";
import { LOGIN_LS_KEY } from "../../constants";

interface Props {
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

const { Header, Content, Sider } = Layout;
const App: React.FC<Props> = ({ setLoggedIn }) => {
  const navigate = useNavigate();

  const userName = localStorage.getItem(LOGIN_LS_KEY);

  const logout = () => {
    localStorage.removeItem(LOGIN_LS_KEY);
    setLoggedIn(false);
    navigate("/login");
  };

  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const items1: MenuProps["items"] = ["1", "2", "3"].map((key) => ({
    key,
    label: `nav ${key}`,
  }));

  const items2: MenuProps["items"] = [
    ...[UserOutlined, LaptopOutlined, NotificationOutlined].map(
      (icon, index) => {
        const key = String(index + 1);

        return {
          key: `sub${key}`,
          icon: React.createElement(icon),
          label: `subnav ${key}`,

          children: new Array(4).fill(null).map((_, j) => {
            const subKey = index * 4 + j + 1;
            return {
              key: subKey,
              label: `option${subKey}`,
            };
          }),
        };
      }
    ),
    {
      key: "logout",
      icon: <LogoutOutlined />,
      label: "Logout",
      onClick: () => {
        logout();
      },
    },
  ];

  return (
    <Layout>
      <Header style={{ display: "flex", alignItems: "center" }}>
        <div className="demo-logo" />
        <Menu
          theme="dark"
          mode="horizontal"
          defaultSelectedKeys={["2"]}
          items={items1}
        />
      </Header>
      <Layout>
        <Sider width={200} style={{ background: colorBgContainer }}>
          <Menu
            mode="inline"
            defaultSelectedKeys={["1"]}
            defaultOpenKeys={["sub1"]}
            style={{ height: "100%", borderRight: 0 }}
            items={items2}
          />
        </Sider>
        <Layout style={{ padding: "0 24px 24px" }}>
          <Breadcrumb style={{ margin: "16px 0" }}>
            <Breadcrumb.Item>Home</Breadcrumb.Item>
            <Breadcrumb.Item>List</Breadcrumb.Item>
            <Breadcrumb.Item>App</Breadcrumb.Item>
          </Breadcrumb>
          <Content
            style={{
              padding: 24,
              margin: 0,
              minHeight: 280,
              background: colorBgContainer,
            }}
          >
            <p style={{ color: "#333", fontWeight: 900, fontSize: "1rem" }}>
              Welcome{" "}
              <span style={{ color: "tomato" }}>{userName?.toUpperCase()}</span>
            </p>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default App;
