package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// openSocks5Conn 建立一个到SOCKS5代理服务器的连接
// 参数:
//
//	socks5Addr: SOCKS5代理服务器的地址，格式为 "host:port"
//	username: SOCKS5认证用的用户名
//	password: SOCKS5认证用的密码
//
// 返回:
//
//	net.Conn: 成功建立的连接，如果连接失败则返回nil
func openSocks5Conn(socks5Addr, username, password string) net.Conn {
	// 尝试建立TCP连接到SOCKS5代理服务器
	conn, err := net.Dial("tcp", socks5Addr)
	if err != nil {
		// 如果连接失败，打印错误信息并返回nil
		fmt.Println("无法连接到SOCKS5服务器:", err)
		return nil
	}

	// 设置连接的读写超时
	// 这有助于防止连接在网络问题时无限期挂起
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// 执行SOCKS5认证握手
	// socks5HandshakeWithAuth函数处理SOCKS5协议的认证过程
	if err := socks5HandshakeWithAuth(conn, username, password); err != nil {
		// 如果握手失败，打印错误信息
		fmt.Println("SOCKS5认证握手失败:", err)
		// 关闭连接以防止资源泄漏
		conn.Close()
		// 返回nil表示连接失败
		return nil
	}

	// 清除之前设置的超时
	// 这允许后续操作不受时
	return conn
}

// socks5HandshakeWithAuth 执行SOCKS5协议的认证握手过程
// 参数:
//
//	conn: 已建立的网络连接
//	username: SOCKS5认证用的用户名
//	password: SOCKS5认证用的密码
//
// 返回:
//
//	error: 如果握手过程中出现错误，返回相应的错误；如果成功，返回nil
func socks5HandshakeWithAuth(conn net.Conn, username, password string) error {
	// 为连接创建一个读取器和写入器
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	// 发送问候消息给SOCKS5服务器
	// 0x05: SOCKS5版本
	// 0x02: 支持的认证方法数量
	// 0x00: 不需要认证
	// 0x02: 用户名/密码认证
	greeting := []byte{0x05, 0x02, 0x00, 0x02}
	if _, err := bw.Write(greeting); err != nil {
		return fmt.Errorf("发送问候消息失败: %v", err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区失败: %v", err)
	}

	// 读取服务器的响应
	response := make([]byte, 2)
	if _, err := br.Read(response); err != nil {
		return fmt.Errorf("读取服务器响应失败: %v", err)
	}

	// 检查服务器是否选择了用户名/密码认证方法
	if response[0] != 0x05 || response[1] != 0x02 {
		fmt.Println("SOCKS5服务器没有选择用户名/密码认证方法")
		return nil
	}

	// 发送认证请求
	authRequest := make([]byte, 3+len(username)+len(password))
	authRequest[0] = 0x01                              // 子协商版本号
	authRequest[1] = byte(len(username))               // 用户名长度
	copy(authRequest[2:], username)                    // 用户名
	authRequest[2+len(username)] = byte(len(password)) // 密码长度
	copy(authRequest[3+len(username):], password)      // 密码
	if _, err := bw.Write(authRequest); err != nil {
		return fmt.Errorf("发送认证请求失败: %v", err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新认证请求缓冲区失败: %v", err)
	}

	// 读取服务器的认证响应
	authResponse := make([]byte, 2)
	if _, err := br.Read(authResponse); err != nil {
		return fmt.Errorf("读取认证响应失败: %v", err)
	}

	// 检查认证是否成功
	// 0x01: 子协商版本号
	// 0x00: 认证成功
	if authResponse[0] != 0x01 || authResponse[1] != 0x00 {
		return fmt.Errorf("SOCKS5认证失败")
	}

	// 认证成功
	return nil
}
