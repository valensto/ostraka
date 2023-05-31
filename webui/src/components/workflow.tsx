import React from 'react';
import {Col, Row, Table as AntdTable, Tag, Typography} from 'antd';
import type {ColumnsType} from 'antd/es/table';

const {Title} = Typography;

interface DataType {
    key: string;
    info: string;
    payload: string;
    status: string;
}

const columns: ColumnsType<DataType> = [
    {
        title: 'Info',
        dataIndex: 'info',
        key: 'info',
    },
    {
        title: 'Payload',
        dataIndex: 'payload',
        key: 'payload',
        render: (text) => <a>{text}</a>,
    },
    {
        title: 'Status',
        key: 'status',
        dataIndex: 'status',
        render: (_, {status}) => {
            let color = 'green';
            if (status === 'failed') {
                color = 'volcano';
            }
            return (
                <Tag color={color} key={status}>
                    {status.toUpperCase()}
                </Tag>
            );
        },
    },
];

const data: DataType[] = [
    {
        key: '1',
        info: 'MQTT',
        payload: "{\"name\": \"John Brown\", \"age\": 32, \"address\": \"New York No. 1 Lake Park\"}",
        status: "dispatched",
    },
    {
        key: '2',
        info: 'Webhook',
        payload: "{\"name\": \"John Brown\", \"age\": 32, \"address\": \"New York No. 1 Lake Park\"}",
        status: "dispatched",
    },
    {
        key: '3',
        info: 'MQTT',
        payload: "{\"name\": \"John Brown\", \"age\": 32, \"address\": \"New York No. 1 Lake Park\"}",
        status: "failed",
    },
];


export const Workflow: React.FC = () => (
    <Row gutter={16}>
        <Col span={12}>
            <Title level={3}>Entries</Title>
            <AntdTable pagination={false} columns={columns} dataSource={data}/>
        </Col>
        <Col span={12}>
            <Title level={3}>Outputs</Title>
            <AntdTable pagination={false} columns={columns} dataSource={data}/>
        </Col>
    </Row>
)