// -----------------------------------------------------------------------
// <copyright file="Ntp.cs" company="Uhuru Software, Inc.">
// Copyright (c) 2011 Uhuru Software, Inc., All Rights Reserved
// </copyright>
// -----------------------------------------------------------------------

namespace BoshUtilities
{
    using System;
    using System.Globalization;
    using System.Net.Sockets;
    using System.Threading;
    using Uhuru.Utilities;
    using System.Diagnostics;
    using System.Text.RegularExpressions;
    using System.Text;
    using Uhuru.BOSH.Agent.Objects;

    /// <summary>
    /// A class the connects to a specified time server and returns the offset
    /// </summary>
    [System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Naming", "CA1704:IdentifiersShouldBeSpelledCorrectly", MessageId = "Ntp", Justification = "FxCop Bug")]
    public class Ntp
    {
        private DateTime currentTime;
        private double offset;
        private string message;

        /// <summary>
        /// Gets the current time.
        /// </summary>
        public DateTime CurrentTime
        {
            get
            {
                return currentTime;
            }
            set
            {
                currentTime = value;
            }
        }

        /// <summary>
        /// Gets the offset.
        /// </summary>
        public double Offset
        {
            get
            {
                return offset;
            }
            set
            {
                offset = value;
            }
        }

        /// <summary>
        /// Gets the connection error message.
        /// </summary>
        public string Message
        {
            get
            {
                return message;
            }
            set
            {
                message = value;
            }
        }

        /// <summary>
        /// Prevents a default instance of the <see cref="Ntp"/> class from being created.
        /// </summary>
        private Ntp()
        {

        }

        /// <summary>
        /// Gets the NTP offset using the default time server.
        /// </summary>
        /// <returns></returns>
        public static NtpMessage GetNtpOffset()
        {
            try
            {
                StringBuilder output = new StringBuilder();
                using (Process w32tm = new Process())
                {
                    ProcessStartInfo info = new ProcessStartInfo();
                    info.Arguments = "/query /status /verbose";
                    info.FileName = "w32tm";
                    info.RedirectStandardOutput = true;
                    info.RedirectStandardInput = true;
                    info.CreateNoWindow = true;
                    info.UseShellExecute = false;
                    w32tm.StartInfo = info;
                    w32tm.EnableRaisingEvents = true;
                    w32tm.OutputDataReceived += new DataReceivedEventHandler(
                        delegate(object sender, DataReceivedEventArgs e)
                        {
                            output.Append(e.Data);
                        }
                        );
                    w32tm.Start(); w32tm.BeginOutputReadLine();
                    w32tm.WaitForExit();
                    w32tm.CancelOutputRead();
                }

                int exitCode = 0;
                if (int.TryParse(Regex.Match(output.ToString(), @"Last Sync Error:\s\d*", RegexOptions.None).Value.Replace("Last Sync Error:", "").Trim(), out exitCode))
                {
                    if (exitCode == 0)
                    {
                        double offset;
                        if (double.TryParse(Regex.Match(output.ToString(), @"Phase Offset:\s\d*.\d*", RegexOptions.None).Value.Replace("Phase Offset:", "").Trim(), out offset))
                        {
                            NtpMessage currentNtp = new NtpMessage();
                            currentNtp.Offset = offset.ToString(CultureInfo.InvariantCulture);
                            currentNtp.Timestamp = DateTime.Now.ToString("dd MMM HH:mm:ss", CultureInfo.InvariantCulture);
                            return currentNtp;
                        }
                    }
                }
                return new NtpMessage() { Message = "bad ntp server" };
            }
            catch (Exception ex)
            {
                Logger.Warning("Could not get NTP offset: {0}", ex.ToString());
                return new NtpMessage() { Message = "bad ntp server" };
            }
        }

        /// <summary>
        /// Gets the NTP offset from a specified time server.
        /// </summary>
        /// <param name="timeServer">The time server.</param>
        /// <returns></returns>
        public static Ntp GetNtpOffset(string timeserver)
        {
            if (string.IsNullOrEmpty(timeserver))
            {
                throw new ArgumentNullException("timeserver");
            }

            Logger.Debug("Retrieving NTP information from {0}", timeserver);

            int retryCount = 5;
            Ntp currentNtp = new Ntp();
            while (retryCount > 0)
            {
                try
                {
                    NtpClient ntpClient = new NtpClient();
                    ntpClient.Init(timeserver);
                    ntpClient.Connect(false);
                    currentNtp.offset = ntpClient.LocalClockOffset;
                    currentNtp.currentTime = DateTime.Now;
                    break;
                }
                catch (SocketException se)
                {
                    Logger.Error("Error while retrieving ntp information: {0}", se.ToString());
                    currentNtp.message = se.Message;
                    retryCount--;
                    Thread.Sleep(1000);
                }
                catch (Exception ex)
                {
                    Logger.Error("Error while retrieving ntp information: {0}", ex.ToString());
                    currentNtp.message = ex.Message;
                    break;
                }
            }
            return currentNtp;
            
        }

        public static void SetTime(double timeOffset)
        {
            Logger.Debug("Updating time");

            Uhuru.BOSH.Agent.NativeMethods.Systemtime st;

            DateTime trts = DateTime.Now.AddMilliseconds(timeOffset);
            st.year = (short)trts.Year;
            st.month = (short)trts.Month;
            st.dayOfWeek = (short)trts.DayOfWeek;
            st.day = (short)trts.Day;
            st.hour = (short)trts.Hour;
            st.minute = (short)trts.Minute;
            st.second = (short)trts.Second;
            st.milliseconds = (short)trts.Millisecond;

            Uhuru.BOSH.Agent.NativeMethods.SetLocalTime(ref st);

            Logger.Debug("Updated local time: {0}", DateTime.Now.ToString(CultureInfo.InvariantCulture));
        }
    }
}
