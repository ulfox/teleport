/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React from 'react';
import AuditEvents from './AuditEvents';
import { makeEvent } from 'teleport/services/audit';

export default {
  title: 'Teleport/Audit',
};

export const Events = () => <AuditEvents {...defaultProps} />;

const events = [
  {
    code: 'T1004I',
    uid: 'b121fc4c-e419-56a2-a760-19cd746c0650',
    time: '2020-06-05T16:24:05Z',
    event: 'user.delete',
    name: 'bob',
    user: 'benarent',
  },
  {
    code: 'T1003I',
    event: 'user.update',
    name: 'bob',
    time: '2020-06-05T16:24:05Z',
    uid: '3a8cd55b5-bce9-5a4c-882d-8e0a5ae10008',
    expires: 111111,
    roles: ['root'],
  },
  {
    code: 'T4002I',
    event: 'session.network',
    namespace: 'default',
    sid: '44c6cea8-362f-11ea-83aa-125400432324',
    server_id: '96f2bed2',
    login: 'root',
    user: 'benarent',
    pid: 2653,
    cgroup_id: 4294968064,
    program: 'bash',
    src_addr: '10.217.136.161',
    dst_addr: '190.58.129.4',
    dst_port: '3000',
    version: 4,
    time: '2019-04-22T19:39:26.676Z',
  },
  {
    code: 'T4001I',
    event: 'session.disk',
    namespace: 'default',
    sid: '44c6cea8-362f-11ea-83aa-125400432324',
    server_id: '96f2bed2',
    login: 'root',
    user: 'benarent',
    pid: 2653,
    cgroup_id: 4294968064,
    program: 'bash',
    path: '/etc/profile.d/',
    flags: 2100000,
    return_code: 0,
    time: '2019-04-22T19:39:26.676Z',
  },
  {
    argv: ['google.com'],
    cgroup_id: 4294968064,
    code: 'T4000I',
    ei: 5,
    event: 'session.command',
    login: 'root',
    namespace: 'default',
    path: '/bin/ping',
    pid: 2653,
    ppid: 2660,
    program: 'ping',
    return_code: 0,
    server_id: '96f2bed2-ebd1-494a-945c-2fd57de41644',
    sid: '44c6cea8-362f-11ea-83aa-125400432324',
    time: '2020-01-13T18:05:53.919Z',
    uid: '734930bb-00e6-4ee6-8798-37f1e9473fac',
    user: 'benarent',
  },
  {
    id: '66b827b2-1b0b-512b-965d-6c789388d3c9',
    code: 'T5000I',
    event: 'access_request.create',
    time: '2020-06-05T19:26:53Z',
    uid: '68a83a99-73ce-4bd7-bbf7-99103c2ba6a0',
    user: 'Carrie_Sandoval',
    state: 'PENDING',
    roles: ['admin'],
  },
  {
    id: '66b827b2-1b0b-512b-965d-6c789388d3c9',
    code: 'T5001I',
    event: 'access_request.update',
    time: '2020-06-05T19:26:53Z',
    uid: '68a83a99-73ce-4bd7-bbf7-99103c2ba6a0',
    state: 'APPROVED',
    updated_by: 'Sam_Waters',
  },
  {
    'addr.local': '172.10.1.1:3022',
    'addr.remote': '172.10.1.254:46992',
    code: 'T2006I',
    ei: 2147483646,
    event: 'session.data',
    login: 'root',
    rx: 3974,
    server_id: 'b331fb6c-85f9-4cb0-b308-3452420bf81e',
    sid: '5fc8bf85-a73e-11ea-afd1-0242ac0a0101',
    time: '2020-06-05T15:14:51Z',
    tx: 4730,
    uid: '2f2f07d0-8a01-4abe-b1c0-5001fd86829b',
    user: 'Stanley_Cooper',
  },
  {
    code: 'T6000I',
    name: 'hello',
    event: 'reset_password_token.create',
    time: '2020-06-05T16:24:22Z',
    ttl: '8h0m0s',
    uid: '85fef5df-6dca-475e-a049-393f4cf1d6a3',
    user: 'b331fb6c-85f9-4cb0-b308-3452420bf81e.one',
  },
  {
    code: 'T8000I',
    event: 'github.created',
    name: 'new_github_connector',
    time: '2020-06-05T19:28:00Z',
    uid: '2b7bb323-35d1-4b9c-9a6d-00ab34c95fb8',
    user: 'unimplemented',
  },
  {
    code: 'T8001I',
    event: 'github.deleted',
    name: 'new_github_connector',
    time: '2020-06-05T19:28:28Z',
    uid: '26f12a67-d593-40df-b3d3-965faee60143',
    user: 'unimplemented',
  },
  {
    code: 'T8100I',
    event: 'oidc.created',
    name: 'new_oidc_connector',
    time: '2020-06-05T19:29:14Z',
    uid: '6208b4b9-0077-41aa-967a-f173b6bcc0d3',
    user: 'unimplemented',
  },
  {
    code: 'T1002I',
    connector: 'local',
    name: 'hello',
    event: 'user.create',
    expires: '0001-01-01T00:00:00Z',
    roles: ['admin'],
    time: '2020-06-05T16:24:05Z',
    uid: '22a273678c-ee78-5ffc-a298-68a841555c98',
    user: 'b331fb6c-85f9-4cb0-b308-3452420bf81e.one',
  },
  {
    code: 'T1005I',
    event: 'user.password_change',
    time: '2020-06-05T19:26:53Z',
    uid: '68a83a99-73ce-4bd7-bbf7-99103c2ba6a0',
    user: 'Ivan_Jordan',
  },
  {
    'addr.local': '172.10.1.1:3022',
    'addr.remote': '172.10.1.254:46992',
    code: 'T2006I',
    ei: 2147483646,
    event: 'session.data',
    login: 'root',
    rx: 3974,
    server_id: 'b331fb6c-85f9-4cb0-b308-3452420bf81e',
    sid: '5fc8bf85-a73e-11ea-afd1-0242ac0a0101',
    time: '2020-06-05T15:14:51Z',
    tx: 4730,
    uid: '2f2f07d0-8a01-4abe-b1c0-5001fd86829b',
    user: 'Betty_Dixon',
  },

  {
    code: 'T2004I',
    ei: 10,
    event: 'session.end',
    namespace: 'default',
    sid: '9febab45-6491-11e9-80a1-427cfde50f5a',
    time: '2019-04-22T00:00:51.543Z',
    uid: '6bf836ee-197c-453e-98e5-31511935f22a',
    user: 'admin@example.com',
  },
  {
    code: 'T1000I',
    event: 'user.login',
    method: 'local',
    success: true,
    time: '2019-04-22T00:49:03Z',
    uid: '173d6b6e-d613-44be-8ff6-f9f893791ef2',
    user: 'admin@example.com',
  },
  {
    code: 'T3007W',
    error:
      'ssh: principal "fsdfdsf" not in the set of valid principals for given certificate: ["root"]',
    event: 'auth',
    success: false,
    time: '2019-04-22T02:09:06Z',
    uid: '036659d6-fdf7-40a4-aa80-74d6ac73b9c0',
    user: 'admin@example.com',
  },
  {
    code: 'T1000W',
    error: 'user(name="fsdfsdf") not found',
    event: 'user.login',
    method: 'local',
    success: false,
    time: '2019-04-22T18:06:32Z',
    uid: '597bf08b-75b2-4dda-a578-e387c5ce9b76',
    user: 'fsdfsdf',
  },
  {
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '151.181.228.114:51454',
    code: 'T2000I',
    ei: 0,
    event: 'session.start',
    login: 'root',
    namespace: 'default',
    server_id: 'de3800ea-69d9-4d72-a108-97e57f8eb393',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    size: '80:25',
    time: '2019-04-22T19:39:26.676Z',
    uid: '84c07a99-856c-419f-9de5-15560451a116',
    user: 'admin@example.com',
  },
  {
    code: 'T2002I',
    ei: 3,
    event: 'resize',
    login: 'root',
    namespace: 'default',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    size: '80:25',
    time: '2019-04-22T19:39:52.432Z',
    uid: '917d8108-3617-4273-ab37-7bbf8e7c1ab9',
    user: 'admin@example.com',
  },
  {
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '151.181.228.114:51752',
    code: 'T2001I',
    ei: 4,
    event: 'session.join',
    login: 'root',
    namespace: 'default',
    server_id: 'de3800ea-69d9-4d72-a108-97e57f8eb393',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    time: '2019-04-22T19:39:52.434Z',
    uid: '13d26190-289b-41d4-af67-c8c8b0617ebe',
    user: 'admin@example.com',
  },
  {
    action: 'download',
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '127.0.0.1:55594',
    code: 'T3004I',
    event: 'scp',
    login: 'root',
    namespace: 'default',
    path: '~/fsdfsdfsdfsdfs',
    time: '2019-04-22T19:41:23Z',
    uid: '183ca6de-c24b-4f67-854f-163c01245fa1',
    user: 'admin@example.com',
  },
  {
    code: 'GE1000I',
    event: 'role.created',
    name: 'admin323232323f',
    time: '2019-04-22T19:42:37Z',
    uid: '19f9248d-0c24-4453-88bb-648fe61feff1',
    user: 'admin@example.com',
  },
  {
    code: 'GE2000I',
    event: 'role.deleted',
    name: 'admin323232323f',
    time: '2019-04-22T19:43:30Z',
    uid: '5244829f-8d7a-4f70-8978-eb3a8f40743a',
    user: 'admin@example.com',
  },
  {
    code: 'GE1000I',
    event: 'role.created',
    name: 'admin323232323',
    time: '2019-04-22T19:44:04Z',
    uid: '6738a844-b49f-4c05-aa84-dcb92b8d8c36',
    user: 'admin@example.com',
  },
  {
    code: 'G1010I',
    event: 'invite.created',
    name: 'fdsfsdf',
    roles: ['@teleadmin'],
    time: '2019-04-22T19:50:48Z',
    uid: 'c5379148-7740-4f76-b0ae-3ddb2438d2a4',
    user: 'admin@example.com',
  },
  {
    code: 'G1003I',
    event: 'logforwarder.created',
    name: 'forwarder1',
    time: '2019-04-22T19:55:18Z',
    uid: '333a9a50-44eb-4f82-9634-7e830d5e3a86',
    user: 'admin@example.com',
  },
  {
    code: 'G2003I',
    event: 'logforwarder.delete',
    name: 'forwarder1',
    time: '2019-04-22T19:55:55Z',
    uid: 'a0e0ba2d-983e-4b2d-a441-3897fa6c9ef2',
    user: 'admin@example.com',
  },

  {
    code: 'T2004I',
    ei: 10,
    event: 'session.end',
    namespace: 'default',
    sid: '9febab45-6491-11e9-80a1-427cfde50f5a',
    time: '2019-04-22T00:00:51.543Z',
    uid: '6bf836ee-197c-453e-98e5-31511935f22a',
    user: 'admin@example.com',
  },
  {
    code: 'T1000I',
    event: 'user.login',
    method: 'local',
    success: true,
    time: '2019-04-22T00:49:03Z',
    uid: '173d6b6e-d613-44be-8ff6-f9f893791ef2',
    user: 'admin@example.com',
  },
  {
    code: 'T3007W',
    error:
      'ssh: principal "fsdfdsf" not in the set of valid principals for given certificate: ["root"]',
    event: 'auth',
    success: false,
    time: '2019-04-22T02:09:06Z',
    uid: '036659d6-fdf7-40a4-aa80-74d6ac73b9c0',
    user: 'admin@example.com',
  },
  {
    code: 'T1000W',
    error: 'user(name="fsdfsdf") not found',
    event: 'user.login',
    method: 'local',
    success: false,
    time: '2019-04-22T18:06:32Z',
    uid: '597bf08b-75b2-4dda-a578-e387c5ce9b76',
    user: 'fsdfsdf',
  },
  {
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '151.181.228.114:51454',
    code: 'T2000I',
    ei: 0,
    event: 'session.start',
    login: 'root',
    namespace: 'default',
    server_id: 'de3800ea-69d9-4d72-a108-97e57f8eb393',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    size: '80:25',
    time: '2019-04-22T19:39:26.676Z',
    uid: '84c07a99-856c-419f-9de5-15560451a116',
    user: 'admin@example.com',
  },
  {
    code: 'T2002I',
    ei: 3,
    event: 'resize',
    login: 'root',
    namespace: 'default',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    size: '80:25',
    time: '2019-04-22T19:39:52.432Z',
    uid: '917d8108-3617-4273-ab37-7bbf8e7c1ab9',
    user: 'admin@example.com',
  },
  {
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '151.181.228.114:51752',
    code: 'T2001I',
    ei: 4,
    event: 'session.join',
    login: 'root',
    namespace: 'default',
    server_id: 'de3800ea-69d9-4d72-a108-97e57f8eb393',
    sid: '56408539-6536-11e9-80a1-427cfde50f5a',
    time: '2019-04-22T19:39:52.434Z',
    uid: '13d26190-289b-41d4-af67-c8c8b0617ebe',
    user: 'admin@example.com',
  },
  {
    code: 'T3004I',
    action: 'download',
    'addr.local': '172.31.28.130:3022',
    'addr.remote': '127.0.0.1:55594',
    event: 'scp',
    login: 'root',
    namespace: 'default',
    path: '~/fsdfsdfsdfsdfs',
    time: '2019-04-22T19:41:23Z',
    uid: '183ca6de-c24b-4f67-854f-163c01245fa1',
    user: 'admin@example.com',
  },
];

const defaultProps = {
  attempt: {},
  attemptActions: {
    do: () => null,
  },
  onFetchLatest: () => Promise.resolve(),
  onFetch: () => Promise.resolve(),
  searchValue: '',
  events: events.map(e => makeEvent(e)),
};
