import socket
import struct
import sys

def eth_addr(a):
    b = "%.2x:%.2x:%.2x:%.2x:%.2x:%.2x" % (a[0], a[1], a[2], a[3], a[4], a[5])
    return b

def main():
    try:
        s = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.ntohs(0x0003))
    except socket.error as msg:
        print('Socket could not be created. Error Code : ' + str(msg[0]) + ' Message ' + msg[1])
        sys.exit()

    print('Packet Sniffer Started...')

    while True:
        packet = s.recvfrom(65565)
        packet = packet[0]

        # Parse Ethernet Header
        eth_length = 14
        eth_header = packet[:eth_length]
        eth = struct.unpack('!6s6sH', eth_header)
        eth_protocol = socket.ntohs(eth[2])

        print('\n--- Ethernet Header ---')
        print('Destination MAC : ' + eth_addr(eth[0]) + ' Source MAC : ' + eth_addr(eth[1]) + ' Protocol : ' + str(eth_protocol))

        # Parse IP Header (if it's an IP packet)
        if eth_protocol == 8: # 8 is IP protocol
            ip_header = packet[eth_length:20 + eth_length]
            iph = struct.unpack('!BBHHHBBH4s4s', ip_header)

            version_ihl = iph[0]
            version = version_ihl >> 4
            ihl = version_ihl & 0xF
            iph_length = ihl * 4

            ttl = iph[5]
            protocol = iph[6]
            s_addr = socket.inet_ntoa(iph[8])
            d_addr = socket.inet_ntoa(iph[9])

            print('--- IP Header ---')
            print('Version : ' + str(version) + ' IP Header Length : ' + str(ihl) + ' TTL : ' + str(ttl) + ' Protocol : ' + str(protocol) + ' Source IP : ' + str(s_addr) + ' Destination IP : ' + str(d_addr))

            # Parse TCP/UDP/ICMP Header
            if protocol == 6: # TCP Protocol
                tcp_header = packet[iph_length + eth_length:iph_length + eth_length + 20]
                tcph = struct.unpack('!HHLLBBHHH', tcp_header)

                source_port = tcph[0]
                dest_port = tcph[1]
                sequence = tcph[2]
                acknowledgement = tcph[3]
                doff_reserved = tcph[4]
                tcph_length = (doff_reserved >> 4) * 4

                print('--- TCP Header ---')
                print('Source Port : ' + str(source_port) + ' Dest Port : ' + str(dest_port) + ' Sequence Number : ' + str(sequence) + ' Acknowledgement : ' + str(acknowledgement) + ' TCP Header Length : ' + str(tcph_length))

            elif protocol == 17: # UDP Protocol
                udp_header = packet[iph_length + eth_length:iph_length + eth_length + 8]
                udph = struct.unpack('!HHHH', udp_header)

                source_port = udph[0]
                dest_port = udph[1]
                length = udph[2]
                checksum = udph[3]

                print('--- UDP Header ---')
                print('Source Port : ' + str(source_port) + ' Dest Port : ' + str(dest_port) + ' Length : ' + str(length) + ' Checksum : ' + str(checksum))

            elif protocol == 1: # ICMP Protocol
                icmp_header = packet[iph_length + eth_length:iph_length + eth_length + 4]
                icmph = struct.unpack('!BBH', icmp_header)

                icmp_type = icmph[0]
                code = icmph[1]
                checksum = icmph[2]

                print('--- ICMP Header ---')
                print('Type : ' + str(icmp_type) + ' Code : ' + str(code) + ' Checksum : ' + str(checksum))

            else:
                print('Other Protocol (IP): ' + str(protocol))
        else:
            print('Other Protocol (Ethernet): ' + str(eth_protocol))

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-21 02:09:29
import socket
import struct
import sys

def get_mac_addr(bytes_addr):
    bytes_str = map('{:02x}'.format, bytes_addr)
    return ':'.join(bytes_str).upper()

def ethernet_frame(data):
    dest_mac, src_mac, proto = struct.unpack('!6s6sH', data[:14])
    return get_mac_addr(dest_mac), get_mac_addr(src_mac), socket.ntohs(proto), data[14:]

def ipv4_packet(data):
    version_header_length = data[0]
    version = version_header_length >> 4
    header_length = (version_header_length & 15) * 4
    ttl, proto, src, target = struct.unpack('!BBH4s4s', data[8:20])
    return version, header_length, ttl, proto, socket.inet_ntoa(src), socket.inet_ntoa(target), data[header_length:]

def icmp_packet(data):
    icmp_type, code, checksum = struct.unpack('!BBH', data[:4])
    return icmp_type, code, checksum, data[4:]

def tcp_segment(data):
    (src_port, dest_port, sequence, acknowledgement, offset_reserved_flags) = struct.unpack('!HHLLH', data[:14])
    offset = (offset_reserved_flags >> 12) * 4
    flag_urg = (offset_reserved_flags & 32) >> 5
    flag_ack = (offset_reserved_flags & 16) >> 4
    flag_psh = (offset_reserved_flags & 8) >> 3
    flag_rst = (offset_reserved_flags & 4) >> 2
    flag_syn = (offset_reserved_flags & 2) >> 1
    flag_fin = (offset_reserved_flags & 1)
    return src_port, dest_port, sequence, acknowledgement, flag_urg, flag_ack, flag_psh, flag_rst, flag_syn, flag_fin, data[offset:]

def udp_segment(data):
    src_port, dest_port, size = struct.unpack('!HHH', data[:8])
    return src_port, dest_port, size, data[8:]

def main():
    if sys.platform.startswith('win'):
        conn = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_IP)
        print("Sniffing on Windows. Only IP packets and above will be parsed.")
    else:
        conn = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.ntohs(0x0003))
        print("Sniffing on Linux/macOS. Ethernet frames will be parsed.")

    print("Starting packet sniffer...")
    print("-" * 60)

    try:
        while True:
            raw_data, addr = conn.recvfrom(65536)

            if sys.platform.startswith('win'):
                eth_proto = 8
                ip_header_data = raw_data
            else:
                dest_mac, src_mac, eth_proto, ip_header_data = ethernet_frame(raw_data)
                print('\nEthernet Frame:')
                print(f'  Destination: {dest_mac}, Source: {src_mac}, Protocol: {eth_proto}')

            if eth_proto == 8:
                (version, header_length, ttl, proto, src_ip, dest_ip, transport_data) = ipv4_packet(ip_header_data)
                print('IPv4 Packet:')
                print(f'  Version: {version}, Header Length: {header_length}, TTL: {ttl}')
                print(f'  Protocol: {proto}, Source: {src_ip}, Target: {dest_ip}')

                if proto == 1:
                    icmp_type, code, checksum, _ = icmp_packet(transport_data)
                    print('ICMP Packet:')
                    print(f'  Type: {icmp_type}, Code: {code}, Checksum: {checksum}')

                elif proto == 6:
                    (src_port, dest_port, sequence, acknowledgement, urg, ack, psh, rst, syn, fin, _) = tcp_segment(transport_data)
                    print('TCP Segment:')
                    print(f'  Source Port: {src_port}, Destination Port: {dest_port}')
                    print(f'  Sequence: {sequence}, Acknowledgment: {acknowledgement}')
                    print(f'  Flags:')
                    print(f'    URG: {urg}, ACK: {ack}, PSH: {psh}, RST: {rst}, SYN: {syn}, FIN: {fin}')

                elif proto == 17:
                    src_port, dest_port, size, _ = udp_segment(transport_data)
                    print('UDP Segment:')
                    print(f'  Source Port: {src_port}, Destination Port: {dest_port}, Length: {size}')

                else:
                    print(f'Other IPv4 Protocol: {proto}')
            else:
                print(f'Other Ethernet Protocol: {eth_proto}')
            print("-" * 60)

    except KeyboardInterrupt:
        print("\nSniffer stopped.")
    except PermissionError:
        print("\nPermission denied. Please run as root/administrator.")
    except Exception as e:
        print(f"\nAn error occurred: {e}")
    finally:
        conn.close()

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-21 02:10:16
import socket
import struct
import sys

def ethernet_frame(data):
    dest_mac, src_mac, proto = struct.unpack('!6s6sH', data[:14])
    return get_mac_addr(dest_mac), get_mac_addr(src_mac), socket.ntohs(proto), data[14:]

def get_mac_addr(bytes_addr):
    bytes_str = map('{:02x}'.format, bytes_addr)
    return ':'.join(bytes_str).upper()

def ipv4_packet(data):
    version_header_length = data[0]
    version = version_header_length >> 4
    header_length = (version_header_length & 15) * 4
    
    (tos, total_length, identification, flags_fragment_offset, ttl, proto, header_checksum, src, target) = \
        struct.unpack('!BH H H BBH 4s 4s', data[1:20])
    
    src_ip = socket.inet_ntoa(src)
    target_ip = socket.inet_ntoa(target)
    
    flags = (flags_fragment_offset >> 13) & 0x7
    fragment_offset = flags_fragment_offset & 0x1FFF

    return version, header_length, tos, total_length, identification, flags, fragment_offset, ttl, proto, header_checksum, src_ip, target_ip, data[header_length:]

def icmp_packet(data):
    icmp_type, code, checksum = struct.unpack('!BBH', data[:4])
    if icmp_type == 8 or icmp_type == 0:
        identifier, sequence_number = struct.unpack('!HH', data[4:8])
        return icmp_type, code, checksum, identifier, sequence_number, data[8:]
    return icmp_type, code, checksum, None, None, data[4:]

def tcp_segment(data):
    (src_port, dest_port, sequence, acknowledgement, offset_reserved_flags, window_size, checksum, urgent_pointer) = \
        struct.unpack('!HHLLHHHH', data[:20])
    offset = (offset_reserved_flags >> 12) * 4
    
    flag_urg = (offset_reserved_flags & 32) >> 5
    flag_ack = (offset_reserved_flags & 16) >> 4
    flag_psh = (offset_reserved_flags & 8) >> 3
    flag_rst = (offset_reserved_flags & 4) >> 2
    flag_syn = (offset_reserved_flags & 2) >> 1
    flag_fin = (offset_reserved_flags & 1)
    
    return src_port, dest_port, sequence, acknowledgement, flag_urg, flag_ack, flag_psh, flag_rst, flag_syn, flag_fin, window_size, checksum, urgent_pointer, data[offset:]

def udp_segment(data):
    src_port, dest_port, size, checksum = struct.unpack('!HHHH', data[:8])
    return src_port, dest_port, size, checksum, data[8:]

def main():
    try:
        conn = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.ntohs(0x0003))
    except socket.error as msg:
        sys.stderr.write('Socket could not be created. Error Code: ' + str(msg[0]) + ' Message: ' + msg[1] + '\n')
        sys.stderr.write('Note: This script requires root/administrator privileges and is primarily designed for Linux/macOS.\n')
        sys.exit(1)

    sys.stdout.write("Starting packet sniffer...\n")
    sys.stdout.write("--------------------------------------------------------------------------------\n")

    while True:
        raw_data, addr = conn.recvfrom(65536)

        dest_mac, src_mac, eth_proto, data = ethernet_frame(raw_data)
        sys.stdout.write('\nEthernet Frame:\n')
        sys.stdout.write(f'  Destination: {dest_mac}, Source: {src_mac}, Protocol: {eth_proto}\n')

        if eth_proto == 8:
            (version, header_length, tos, total_length, identification, flags, fragment_offset, ttl, proto, header_checksum, src_ip, target_ip, data) = ipv4_packet(data)
            sys.stdout.write('  IPv4 Packet:\n')
            sys.stdout.write(f'    Version: {version}, Header Length: {header_length}, Type of Service: {tos}\n')
            sys.stdout.write(f'    Total Length: {total_length}, ID: {identification}, Flags: {flags}, Fragment Offset: {fragment_offset}\n')
            sys.stdout.write(f'    TTL: {ttl}, Protocol: {proto}, Header Checksum: {header_checksum}\n')
            sys.stdout.write(f'    Source: {src_ip}, Target: {target_ip}\n')

            if proto == 1:
                icmp_type, code, checksum, identifier, sequence_number, data = icmp_packet(data)
                sys.stdout.write('    ICMP Packet:\n')
                sys.stdout.write(f'      Type: {icmp_type}, Code: {code}, Checksum: {checksum}\n')
                if identifier is not None:
                    sys.stdout.write(f'      Identifier: {identifier}, Sequence Number: {sequence_number}\n')

            elif proto == 6:
                (src_port, dest_port, sequence, acknowledgement, flag_urg, flag_ack, flag_psh, flag_rst, flag_syn, flag_fin, window_size, checksum, urgent_pointer, data) = tcp_segment(data)
                sys.stdout.write('    TCP Segment:\n')
                sys.stdout.write(f'      Source Port: {src_port}, Destination Port: {dest_port}\n')
                sys.stdout.write(f'      Sequence: {sequence}, Acknowledgement: {acknowledgement}\n')
                sys.stdout.write(f'      Flags: URG:{flag_urg}, ACK:{flag_ack}, PSH:{flag_psh}, RST:{flag_rst}, SYN:{flag_syn}, FIN:{flag_fin}\n')
                sys.stdout.write(f'      Window Size: {window_size}, Checksum: {checksum}, Urgent Pointer: {urgent_pointer}\n')

            elif proto == 17:
                src_port, dest_port, size, checksum, data = udp_segment(data)
                sys.stdout.write('    UDP Segment:\n')
                sys.stdout.write(f'      Source Port: {src_port}, Destination Port: {dest_port}, Length: {size}, Checksum: {checksum}\n')

            else:
                sys.stdout.write(f'    Other IPv4 Protocol: {proto}\n')
        
        elif eth_proto == 1544:
            sys.stdout.write('  ARP Packet (not parsed further)\n')

        else:
            sys.stdout.write(f'  Other Ethernet Protocol: {eth_proto}\n')

if __name__ == '__main__':
    main()

# Additional implementation at 2025-06-21 02:11:42
import socket
import struct
import binascii
import sys

def eth_header(data):
    dest_mac, src_mac, eth_type = struct.unpack('!6s6sH', data[:14])
    return (
        binascii.hexlify(dest_mac).decode('utf-8'),
        binascii.hexlify(src_mac).decode('utf-8'),
        socket.ntohs(eth_type),
        data[14:]
    )

def ipv4_header(data):
    version_ihl, dscp_ecn, total_length, identification, flags_fragment_offset, ttl, protocol, hdr_checksum, src_ip, dest_ip = struct.unpack('!BBHHHBBH4s4s', data[:20])

    version = (version_ihl >> 4) & 0xF
    ihl = (version_ihl & 0xF) * 4
    
    src_ip_str = socket.inet_ntoa(src_ip)
    dest_ip_str = socket.inet_ntoa(dest_ip)

    return (
        version,
        ihl,
        total_length,
        protocol,
        src_ip_str,
        dest_ip_str,
        data[ihl:]
    )

def tcp_header(data):
    src_port, dest_port, sequence, acknowledgement, offset_reserved_flags = struct.unpack('!HHLLH', data[:14])

    offset = ((offset_reserved_flags >> 12) & 0xF) * 4

    window_size, checksum, urgent_pointer = struct.unpack('!HHH', data[14:20])

    return (
        src_port,
        dest_port,
        sequence,
        acknowledgement,
        offset,
        window_size,
        checksum,
        urgent_pointer,
        data[offset:]
    )

def udp_header(data):
    src_port, dest_port, length, checksum = struct.unpack('!HHHH', data[:8])
    return (
        src_port,
        dest_port,
        length,
        checksum,
        data[8:]
    )

def icmp_header(data):
    icmp_type, code, checksum = struct.unpack('!BBH', data[:4])
    return (
        icmp_type,
        code,
        checksum,
        data[4:]
    )

def sniffer(host):
    try:
        s = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.ntohs(0x0003))
    except AttributeError:
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_IP)
            if sys.platform == 'win32':
                s.ioctl(socket.SIO_RCVALL, socket.RCVALL_ON)
        except socket.error as e:
            print(f"Error creating socket: {e}")
            sys.exit(1)
    except socket.error as e:
        print(f"Error creating socket: {e}")
        sys.exit(1)

    print("Starting packet sniffer...")
    print("--------------------------------------------------------------------")

    try:
        while True:
            raw_data, addr = s.recvfrom(65535)

            if s.family == socket.AF_PACKET:
                eth_dest_mac, eth_src_mac, eth_type, ip_data = eth_header(raw_data)
                print(f"\n[Ethernet Frame]")
                print(f"  Destination MAC: {eth_dest_mac}")
                print(f"  Source MAC:      {eth_src_mac}")
                print(f"  EtherType:       0x{eth_type:04x}")

                if eth_type == 0x0800:
                    ip_version, ip_ihl, ip_total_length, ip_protocol, ip_src_ip, ip_dest_ip, transport_data = ipv4_header(ip_data)
                    print(f"[IPv4 Packet]")
                    print(f"  Version:         {ip_version}")
                    print(f"  Header Length:   {ip_ihl} bytes")
                    print(f"  Total Length:    {ip_total_length} bytes")
                    print(f"  Source IP:       {ip_src_ip}")
                    print(f"  Destination IP:  {ip_dest_ip}")
                    print(f"  Protocol:        {ip_protocol}")

                    if ip_protocol == 6:
                        tcp_src_port, tcp_dest_port, tcp_sequence, tcp_acknowledgement, tcp_offset, tcp_window, tcp_checksum, tcp_urgent, payload = tcp_header(transport_data)
                        print(f"[TCP Segment]")
                        print(f"  Source Port:     {tcp_src_port}")
                        print(f"  Destination Port:{tcp_dest_port}")
                        print(f"  Sequence:        {tcp_sequence}")
                        print(f"  Acknowledgement: {tcp_acknowledgement}")
                        print(f"  Header Length:   {tcp_offset} bytes")
                        print(f"  Window Size:     {tcp_window}")
                        print(f"  Checksum:        0x{tcp_checksum:04x}")
                        print(f"  Urgent Pointer:  {tcp_urgent}")
                    elif ip_protocol == 17:
                        udp_src_port, udp_dest_port, udp_length, udp_checksum, payload = udp_header(transport_data)
                        print(f"[UDP Datagram]")
                        print(f"  Source Port:     {udp_src_port}")
                        print(f"  Destination Port:{udp_dest_port}")
                        print(f"  Length:          {udp_length} bytes")
                        print(f"  Checksum:        0x{udp_checksum:04x}")
                    elif ip_protocol == 1:
                        icmp_type, icmp_code, icmp_checksum, payload = icmp_header(transport_data)
                        print(f"[ICMP Message]")
                        print(f"  Type:            {icmp_type}")
                        print(f"  Code:            {icmp_code}")
                        print(f"  Checksum:        0x{icmp_checksum:04x}")
                    else:
                        print(f"[Unknown IP Protocol: {ip_protocol}]")
                elif eth_type == 0x0806:
                    print("[ARP Packet detected (not parsed)]")
                else:
                    print(f"[Unknown EtherType: 0x{eth_type:04x}]")

            elif s.family == socket.AF_INET:
                ip_version, ip_ihl, ip_total_length, ip_protocol, ip_src_ip, ip_dest_ip, transport_data = ipv4_header(raw_data)
                print(f"\n[IPv4 Packet (Windows)]")
                print(f"  Version:         {ip_version}")
                print(f"  Header Length:   {ip_ihl} bytes")
                print(f"  Total Length:    {ip_total_length} bytes")
                print(f"  Source IP:       {ip_src_ip}")
                print(f"  Destination IP:  {ip_dest_ip}")
                print(f"  Protocol:        {ip_protocol}")

                if ip_protocol == 6:
                    tcp_src_port, tcp_dest_port, tcp_sequence, tcp_acknowledgement, tcp_offset, tcp_window, tcp_checksum, tcp_urgent, payload = tcp_header(transport_data)
                    print(f"[TCP Segment]")
                    print(f"  Source Port:     {tcp_src_port}")
                    print(f"  Destination Port:{tcp_dest_port}")
                    print(f"  Sequence:        {tcp_sequence}")
                    print(f"  Acknowledgement: {tcp_acknowledgement}")
                    print(f"  Header Length:   {tcp_offset} bytes")
                    print(f"  Window Size:     {tcp_window}")
                    print(f"  Checksum:        0x{tcp_checksum:04x}")
                    print(f"  Urgent Pointer:  {tcp_urgent}")
                elif ip_protocol == 17:
                    udp_src_port, udp_dest_port, udp_length, udp_checksum, payload = udp_header(transport_data)
                    print(f"[UDP Datagram]")
                    print(f"  Source Port:     {udp_src_port}")
                    print(f"  Destination Port:{udp_dest_port}")
                    print(f"  Length:          {udp_length} bytes")
                    print(f"  Checksum:        0x{udp_checksum:04x}")
                elif ip_protocol == 1:
                    icmp_type, icmp_code, icmp_checksum, payload = icmp_header(transport_data)
                    print(f"[ICMP Message]")
                    print(f"  Type:            {icmp_type}")
                    print(f"  Code:            {icmp_code}")
                    print(f"  Checksum:        0x{icmp_checksum:04x}")
                else:
                    print(f"[Unknown IP Protocol: {ip_protocol}]")

            print("--------------------------------------------------------------------")

    except KeyboardInterrupt:
        print("\nSniffer stopped.")
    finally:
        if s.family == socket.AF_INET and sys.platform == 'win32':
            s.ioctl(socket.SIO_RCVALL, socket.RCVALL_OFF)
        s.close()

if __name__ == '__main__':
    sniffer('0.0.0.0')