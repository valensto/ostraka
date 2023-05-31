import React from 'react';
import {Layout, Menu, Tabs, TabsProps} from 'antd';
import {Workflow} from "./components/workflow.tsx";

const { Header, Content, Footer } = Layout;

const onChange = (key: string) => {
    console.log(key);
};

const items: TabsProps['items'] = [
    {
        key: '1',
        label: `Orders`,
        children: <Workflow />,
    },
    {
        key: '2',
        label: `Accounts`,
        children: <Workflow />,
    },
    {
        key: '3',
        label: `Feedback`,
        children: <Workflow />,
    },
];

const App: React.FC = () => {
    return (
        <Layout className="layout">
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <h3 style={{color: "white"}}>Ostraka</h3>
                <Menu
                    theme="dark"
                    mode="horizontal"
                    items={[
                        {
                            key: '1',
                            label: 'Documentation',
                        }
                    ]}
                />
            </Header>

            <Content style={{ padding: '0 50px', minHeight: "90vh" }}>
                <Tabs defaultActiveKey="1" items={items} onChange={onChange} />
            </Content>
            <Footer style={{ textAlign: 'center' }}>Created by Ostraka ©2023</Footer>
        </Layout>
    );
};

export default App;