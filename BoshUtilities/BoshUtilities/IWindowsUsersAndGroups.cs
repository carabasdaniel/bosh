using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;

namespace BoshUtilities
{
    [ComVisible(true)]
    public interface IWindowsUsersAndGroups
    {
        string[] GetUsers();
        Dictionary<string, string> GetUsersDescription();
        void CreateUser(string userName, string password, string basepath);
        void SetUserPassword(string userName, string password);
        void DeleteUser(string userName);
        bool ExistsUser(string userName);
        string[] GetGroups();
        void CreateGroup(string groupName);
        void DeleteGroup(string groupName);
        bool ExistsGroup(string groupName);
        void AddUserToGroup(string userName, string groupName);
        void RemoveUserFromGroup(string userName, string groupName);
        bool IsUserMemberOfGroup(string userName, string groupName);
        string GetLocalUserSid(string userName);
    }
}
