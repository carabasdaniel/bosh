using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace BoshUtilities
{
    using System;
    using System.Net;
    using System.Net.Sockets;
    using System.Runtime.InteropServices;
    using System.Globalization;

    // Leap indicator field values
    public enum LeapIndicator
    {
        NoWarning,		// 0 - No warning
        LastMinute61,	// 1 - Last minute has 61 seconds
        LastMinute59,	// 2 - Last minute has 59 seconds
        Alarm			// 3 - Alarm condition (clock not synchronized)
    }

    //Mode field values
    public enum Mode
    {
        SymmetricActive,	// 1 - Symmetric active
        SymmetricPassive,	// 2 - Symmetric pasive
        Client,				// 3 - Client
        Server,				// 4 - Server
        Broadcast,			// 5 - Broadcast
        Unknown				// 0, 6, 7 - Reserved
    }


    /// <summary>
    /// Stratum field values
    /// </summary>
    public enum Stratum
    {
        /// <summary>
        /// 0 - unspecified or unavailable
        /// </summary>
        Unspecified,

        /// <summary>
        /// 1 - primary reference (e.g. radio-clock)
        /// </summary>
        PrimaryReference,


        /// <summary>
        /// 2-15 - secondary reference (via NTP or SNTP)
        /// </summary>
        SecondaryReference,

        /// <summary>
        /// 16-255 - reserved
        /// </summary>
        [System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Naming", "CA1700:DoNotNameEnumValuesReserved")]
        Reserved				
    }

    /// <summary>
    /// NTPClient is a C# class designed to connect to time servers on the Internet.
    /// The implementation of the protocol is based on the RFC 2030.
    /// 
    /// Public class members:
    ///
    /// LeapIndicator - Warns of an impending leap second to be inserted/deleted in the last
    /// minute of the current day. (See the _LeapIndicator enum)
    /// 
    /// VersionNumber - Version number of the protocol (3 or 4).
    /// 
    /// Mode - Returns mode. (See the _Mode enum)
    /// 
    /// Stratum - Stratum of the clock. (See the _Stratum enum)
    /// 
    /// PollInterval - Maximum interval between successive messages.
    /// 
    /// Precision - Precision of the clock.
    /// 
    /// RootDelay - Round trip time to the primary reference source.
    /// 
    /// RootDispersion - Nominal error relative to the primary reference source.
    /// 
    /// ReferenceID - Reference identifier (either a 4 character string or an IP address).
    /// 
    /// ReferenceTimestamp - The time at which the clock was last set or corrected.
    /// 
    /// OriginateTimestamp - The time at which the request departed the client for the server.
    /// 
    /// ReceiveTimestamp - The time at which the request arrived at the server.
    /// 
    /// Transmit Timestamp - The time at which the reply departed the server for client.
    /// 
    /// RoundTripDelay - The time between the departure of request and arrival of reply.
    /// 
    /// LocalClockOffset - The offset of the local clock relative to the primary reference
    /// source.
    /// 
    /// Initialize - Sets up data structure and prepares for connection.
    /// 
    /// Connect - Connects to the time server and populates the data structure.
    ///	It can also set the system time.
    /// 
    /// IsResponseValid - Returns true if received data is valid and if comes from
    /// a NTP-compliant time server.
    /// 
    /// ToString - Returns a string representation of the object.
    /// 
    /// -----------------------------------------------------------------------------
    /// Structure of the standard NTP header (as described in RFC 2030)
    ///                       1                   2                   3
    ///   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |LI | VN  |Mode |    Stratum    |     Poll      |   Precision   |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                          Root Delay                           |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                       Root Dispersion                         |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                     Reference Identifier                      |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                                                               |
    ///  |                   Reference Timestamp (64)                    |
    ///  |                                                               |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                                                               |
    ///  |                   Originate Timestamp (64)                    |
    ///  |                                                               |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                                                               |
    ///  |                    Receive Timestamp (64)                     |
    ///  |                                                               |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                                                               |
    ///  |                    Transmit Timestamp (64)                    |
    ///  |                                                               |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                 Key Identifier (optional) (32)                |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ///  |                                                               |
    ///  |                                                               |
    ///  |                 Message Digest (optional) (128)               |
    ///  |                                                               |
    ///  |                                                               |
    ///  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    /// 
    /// -----------------------------------------------------------------------------
    /// 
    /// NTP Timestamp Format (as described in RFC 2030)
    ///                         1                   2                   3
    ///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    /// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    /// |                           Seconds                             |
    /// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    /// |                  Seconds Fraction (0-padded)                  |
    /// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    /// 
    /// </summary>

    [ComVisible(true)]
    public class NtpClient:INTPClient
    {
        // NTP Data Structure Length
        private const byte NTPDataLength = 48;
        // NTP Data Structure (as described in RFC 2030)
        byte[] NTPData = new byte[NTPDataLength];

        // Offset constants for timestamps in the data structure
        private const byte offReferenceID = 12;
        private const byte offReferenceTimestamp = 16;
        private const byte offOriginateTimestamp = 24;
        private const byte offReceiveTimestamp = 32;
        private const byte offTransmitTimestamp = 40;

        // Leap Indicator
        public LeapIndicator LeapIndicator
        {
            get
            {
                // Isolate the two most significant bits
                byte val = (byte)(NTPData[0] >> 6);
                switch (val)
                {
                    case 0: return LeapIndicator.NoWarning;
                    case 1: return LeapIndicator.LastMinute61;
                    case 2: return LeapIndicator.LastMinute59;
                    case 3: goto default;
                    default:
                        return LeapIndicator.Alarm;
                }
            }
        }

        // Version Number
        public byte VersionNumber
        {
            get
            {
                // Isolate bits 3 - 5
                byte val = (byte)((NTPData[0] & 0x38) >> 3);
                return val;
            }
        }

        // Mode
        public Mode Mode
        {
            get
            {
                // Isolate bits 0 - 3
                byte val = (byte)(NTPData[0] & 0x7);
                switch (val)
                {
                    case 0: goto default;
                    case 6: goto default;
                    case 7: goto default;
                    default:
                        return Mode.Unknown;
                    case 1:
                        return Mode.SymmetricActive;
                    case 2:
                        return Mode.SymmetricPassive;
                    case 3:
                        return Mode.Client;
                    case 4:
                        return Mode.Server;
                    case 5:
                        return Mode.Broadcast;
                }
            }
        }

        // Stratum
        public Stratum Stratum
        {
            get
            {
                byte val = (byte)NTPData[1];
                if (val == 0) return Stratum.Unspecified;
                else
                    if (val == 1) return Stratum.PrimaryReference;
                    else
                        if (val <= 15) return Stratum.SecondaryReference;
                        else
                            return Stratum.Reserved;
            }
        }

        // Poll Interval
        public uint PollInterval
        {
            get
            {
                return (uint)Math.Round(Math.Pow(2, (sbyte)NTPData[2]));
            }
        }

        // Precision (in milliseconds)
        public double Precision
        {
            get
            {
                return (1000 * Math.Pow(2, (sbyte)NTPData[3]));
            }
        }

        // Root Delay (in milliseconds)
        public double RootDelay
        {
            get
            {
                int temp = 0;
                temp = 256 * (256 * (256 * NTPData[4] + NTPData[5]) + NTPData[6]) + NTPData[7];
                return 1000 * (((double)temp) / 0x10000);
            }
        }

        // Root Dispersion (in milliseconds)
        public double RootDispersion
        {
            get
            {
                int temp = 0;
                temp = 256 * (256 * (256 * NTPData[8] + NTPData[9]) + NTPData[10]) + NTPData[11];
                return 1000 * (((double)temp) / 0x10000);
            }
        }

        /// <summary>
        /// Gets the reference id.
        /// </summary>
        public string ReferenceId
        {
            get
            {
                string val = "";
                switch (Stratum)
                {
                    case Stratum.Unspecified:
                        goto case Stratum.PrimaryReference;
                    case Stratum.PrimaryReference:
                        val += (char)NTPData[offReferenceID + 0];
                        val += (char)NTPData[offReferenceID + 1];
                        val += (char)NTPData[offReferenceID + 2];
                        val += (char)NTPData[offReferenceID + 3];
                        break;
                    case Stratum.SecondaryReference:
                        switch (VersionNumber)
                        {
                            case 3:	// Version 3, Reference ID is an IPv4 address
                                string Address = NTPData[offReferenceID + 0].ToString(CultureInfo.InvariantCulture) + "." +
                                                 NTPData[offReferenceID + 1].ToString(CultureInfo.InvariantCulture) + "." +
                                                 NTPData[offReferenceID + 2].ToString(CultureInfo.InvariantCulture) + "." +
                                                 NTPData[offReferenceID + 3].ToString(CultureInfo.InvariantCulture);
                                try
                                {
                                    IPHostEntry Host = Dns.GetHostByAddress(Address);
                                    val = Host.HostName + " (" + Address + ")";
                                }
                                catch (Exception)
                                {
                                    val = "N/A";
                                }
                                break;
                            case 4: // Version 4, Reference ID is the timestamp of last update
                                DateTime time = ComputeDate(GetMilliSeconds(offReferenceID));
                                // Take care of the time zone
                                TimeSpan offspan = TimeZone.CurrentTimeZone.GetUtcOffset(DateTime.Now);
                                val = (time + offspan).ToString();
                                break;
                            default:
                                val = "N/A";
                                break;
                        }
                        break;
                }

                return val;
            }
        }

        // Reference Timestamp
        public DateTime ReferenceTimestamp
        {
            get
            {
                DateTime time = ComputeDate(GetMilliSeconds(offReferenceTimestamp));
                // Take care of the time zone
                TimeSpan offspan = TimeZone.CurrentTimeZone.GetUtcOffset(DateTime.Now);
                return time + offspan;
            }
        }

        // Originate Timestamp
        public DateTime OriginateTimestamp
        {
            get
            {
                return ComputeDate(GetMilliSeconds(offOriginateTimestamp));
            }
        }

        // Receive Timestamp
        public DateTime ReceiveTimestamp
        {
            get
            {
                DateTime time = ComputeDate(GetMilliSeconds(offReceiveTimestamp));
                // Take care of the time zone
                TimeSpan offspan = TimeZone.CurrentTimeZone.GetUtcOffset(DateTime.Now);
                return time + offspan;
            }
        }

        // Transmit Timestamp
        public DateTime TransmitTimestamp
        {
            get
            {
                DateTime time = ComputeDate(GetMilliSeconds(offTransmitTimestamp));
                // Take care of the time zone
                TimeSpan offspan = TimeZone.CurrentTimeZone.GetUtcOffset(DateTime.Now);
                return time + offspan;
            }
            set
            {
                SetDate(offTransmitTimestamp, value);
            }
        }

        // Reception Timestamp
        private DateTime ReceptionTimestamp;

       
        /// <summary>
        /// Gets the round trip delay (in milliseconds)
        /// </summary>
        public int RoundtripDelay
        {
            get
            {
                TimeSpan span = (ReceiveTimestamp - OriginateTimestamp) + (ReceptionTimestamp - TransmitTimestamp);
                return (int)span.TotalMilliseconds;
            }
        }

        // Local clock offset (in milliseconds)
        public double LocalClockOffset
        {
            get
            {
                TimeSpan span = (ReceiveTimestamp - OriginateTimestamp) - (ReceptionTimestamp - TransmitTimestamp);
                return (double)(span.TotalMilliseconds / 2);
            }
        }

        // Compute date, given the number of milliseconds since January 1, 1900
        private static DateTime ComputeDate(ulong milliseconds)
        {
            TimeSpan span = TimeSpan.FromMilliseconds((double)milliseconds);
            DateTime time = new DateTime(1900, 1, 1);
            time += span;
            return time;
        }

        // Compute the number of milliseconds, given the offset of a 8-byte array
        private ulong GetMilliSeconds(byte offset)
        {
            ulong intpart = 0, fractpart = 0;

            for (int i = 0; i <= 3; i++)
            {
                intpart = 256 * intpart + NTPData[offset + i];
            }
            for (int i = 4; i <= 7; i++)
            {
                fractpart = 256 * fractpart + NTPData[offset + i];
            }
            ulong milliseconds = intpart * 1000 + (fractpart * 1000) / 0x100000000L;
            return milliseconds;
        }

        // Compute the 8-byte array, given the date
        private void SetDate(byte offset, DateTime date)
        {
            ulong intpart = 0, fractpart = 0;
            DateTime StartOfCentury = new DateTime(1900, 1, 1, 0, 0, 0);	// January 1, 1900 12:00 AM

            ulong milliseconds = (ulong)(date - StartOfCentury).TotalMilliseconds;
            intpart = milliseconds / 1000;
            fractpart = ((milliseconds % 1000) * 0x100000000L) / 1000;

            ulong temp = intpart;
            for (int i = 3; i >= 0; i--)
            {
                NTPData[offset + i] = (byte)(temp % 256);
                temp = temp / 256;
            }

            temp = fractpart;
            for (int i = 7; i >= 4; i--)
            {
                NTPData[offset + i] = (byte)(temp % 256);
                temp = temp / 256;
            }
        }

        // Initialize the NTPClient data
        private void Initialize()
        {
            // Set version number to 4 and Mode to 3 (client)
            NTPData[0] = 0x1B;
            // Initialize all other fields with 0
            for (int i = 1; i < 48; i++)
            {
                NTPData[i] = 0;
            }
            // Initialize the transmit timestamp
            TransmitTimestamp = DateTime.Now;
        }

        public NtpClient()
        {

        }      

        [ComVisible(true)]
        public void Init(string host)
        {
            TimeServer = host;
        }

        /// <summary>
        /// Connect to the time server and update system time
        /// </summary>
        /// <param name="updateSystemTime">if set to <c>true</c> [update system time].</param>
        [ComVisible(true)]
        public void Connect(bool updateSystemTime)
        {
            // Resolve server address
            IPHostEntry hostadd = Dns.GetHostEntry(TimeServer);
            IPEndPoint EPhost = new IPEndPoint(hostadd.AddressList[0], 123);

            //Connect the time server
            using (UdpClient timeSocket = new UdpClient())
            {
                timeSocket.Client.SendTimeout = 5000;
                timeSocket.Client.ReceiveTimeout = 5000;
                timeSocket.Connect(EPhost);
                
                // Initialize data structure
                Initialize();
                timeSocket.Send(NTPData, NTPData.Length);
                
                NTPData = timeSocket.Receive(ref EPhost);
                if (!IsResponseValid())
                {
                    throw new InvalidOperationException("Invalid response from " + TimeServer);
                }
                ReceptionTimestamp = DateTime.Now;

            }

            // Update system time
            if (updateSystemTime)
            {
                SetTime();
            }
        }

        // Check if the response from server is valid
        public bool IsResponseValid()
        {
            if (NTPData.Length < NTPDataLength || Mode != Mode.Server)
            {
                return false;
            }
            else
            {
                return true;
            }
        }

        // Converts the object to string

        [ComVisible(true)]
        public override string ToString()
        {
            string str;

            str = "Leap Indicator: ";
            switch (LeapIndicator)
            {
                case LeapIndicator.NoWarning:
                    str += "No warning";
                    break;
                case LeapIndicator.LastMinute61:
                    str += "Last minute has 61 seconds";
                    break;
                case LeapIndicator.LastMinute59:
                    str += "Last minute has 59 seconds";
                    break;
                case LeapIndicator.Alarm:
                    str += "Alarm Condition (clock not synchronized)";
                    break;
            }
            str += "\r\nVersion number: " + VersionNumber.ToString(CultureInfo.InvariantCulture) + "\r\n";
            str += "Mode: ";
            switch (Mode)
            {
                case Mode.Unknown:
                    str += "Unknown";
                    break;
                case Mode.SymmetricActive:
                    str += "Symmetric Active";
                    break;
                case Mode.SymmetricPassive:
                    str += "Symmetric Pasive";
                    break;
                case Mode.Client:
                    str += "Client";
                    break;
                case Mode.Server:
                    str += "Server";
                    break;
                case Mode.Broadcast:
                    str += "Broadcast";
                    break;
            }
            str += "\r\nStratum: ";
            switch (Stratum)
            {
                case Stratum.Unspecified:
                case Stratum.Reserved:
                    str += "Unspecified";
                    break;
                case Stratum.PrimaryReference:
                    str += "Primary Reference";
                    break;
                case Stratum.SecondaryReference:
                    str += "Secondary Reference";
                    break;
            }
            str += "\r\nLocal time: " + TransmitTimestamp.ToString(CultureInfo.InvariantCulture);
            str += "\r\nPrecision: " + Precision.ToString(CultureInfo.InvariantCulture) + " ms";
            str += "\r\nPoll Interval: " + PollInterval.ToString(CultureInfo.InvariantCulture) + " s";
            str += "\r\nReference ID: " + ReferenceId.ToString(CultureInfo.InvariantCulture);
            str += "\r\nRoot Delay: " + RootDelay.ToString(CultureInfo.InvariantCulture) + " ms";
            str += "\r\nRoot Dispersion: " + RootDispersion.ToString(CultureInfo.InvariantCulture) + " ms";
            str += "\r\nRound Trip Delay: " + RoundtripDelay.ToString(CultureInfo.InvariantCulture) + " ms";
            str += "\r\nLocal Clock Offset: " + LocalClockOffset.ToString(CultureInfo.InvariantCulture) + " ms";
            str += "\r\n";

            return str;
        }

       

        // Set system time according to transmit timestamp
        private void SetTime()
        {
            Uhuru.BOSH.Agent.NativeMethods.Systemtime st;

            DateTime trts = DateTime.Now.AddMilliseconds(LocalClockOffset);
            st.year = (short)trts.Year;
            st.month = (short)trts.Month;
            st.dayOfWeek = (short)trts.DayOfWeek;
            st.day = (short)trts.Day;
            st.hour = (short)trts.Hour;
            st.minute = (short)trts.Minute;
            st.second = (short)trts.Second;
            st.milliseconds = (short)trts.Millisecond;

            Uhuru.BOSH.Agent.NativeMethods.SetLocalTime(ref st);
        }

        // The URL of the time server we're connecting to
        private string TimeServer;
    }
}
