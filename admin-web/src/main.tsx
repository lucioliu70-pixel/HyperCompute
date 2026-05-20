import React from 'react';
import ReactDOM from 'react-dom/client';
import { Button, Layout, Tabs, Table } from 'antd';

const token = 'admin-token';
const fetcher = (u: string, opt: any = {}) => fetch(u, { ...opt, headers: { Authorization: `Bearer ${token}`, 'content-type': 'application/json' } }).then(r => r.json());

function App() {
  const [contributors, setContributors] = React.useState<any[]>([]);
  const [nodes, setNodes] = React.useState<any[]>([]);
  React.useEffect(() => { fetcher('/api/admin/contributors').then(setContributors); fetcher('/api/admin/nodes').then(setNodes); }, []);
  return <Layout><Layout.Content style={{ padding: 24 }}>
    <Tabs items={[{ key: 'c', label: '贡献者管理', children: <Table rowKey='user_id' dataSource={contributors} columns={[{title:'user_id',dataIndex:'user_id'},{title:'role',dataIndex:'role'},{title:'status',dataIndex:'contributor_status'},{title:'节点数',dataIndex:'node_count'},{title:'在线',dataIndex:'online_nodes'},{title:'收益',dataIndex:'earnings'},{title:'积分',dataIndex:'points'},{title:'操作',render:(_,r)=><Button onClick={()=>fetcher(`/api/admin/contributors/${r.user_id}/approve`,{method:'POST'}).then(()=>location.reload())}>Approve</Button>}]} /> },
      { key: 'n', label: '节点管理', children: <Table rowKey='node_id' dataSource={nodes} columns={[{title:'node_id',dataIndex:'node_id'},{title:'owner',dataIndex:'owner_user_id'},{title:'base_url',dataIndex:'base_url'},{title:'status',dataIndex:'status'},{title:'runtime_online',dataIndex:'runtime_online'},{title:'gpu_usage',dataIndex:'gpu_usage'},{title:'last_heartbeat_at',dataIndex:'last_heartbeat_at'}]} /> },
      { key: 'p', label: '贡献者门户', children: <pre>步骤: 1申请贡献者 2审批 3运行 node-daemon
NODE_ID=node-001 OWNER_USER_ID=2 SCHEDULER_BASE_URL=http://localhost:8081 go run ./node-daemon</pre> }]} />
  </Layout.Content></Layout>
}
ReactDOM.createRoot(document.getElementById('root')!).render(<App />);
