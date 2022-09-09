using System;
using System.Text;
using System.Runtime.InteropServices;

namespace clientAgent
{
    class clientAgent
    {
        [DllImport("librpc.so.1.0")]
        static extern int rpcExec([In] byte[] rpccmd, ref IntPtr output);

        static void Main(string[] args)
        {
            var res = "";
            foreach (var arg in args)
            {
                res += $"{arg} ";
            }

            // Example commands to be passed in
            // string res = "activate -u wss://192.168.1.96/activate -n -profile Test_Profile";
            // string res = "amtinfo";

            IntPtr output = IntPtr.Zero;
            int returnCode = rpcExec(Encoding.ASCII.GetBytes(res), ref output);
            Console.WriteLine("rpcExec completed: return code[" + returnCode + "] " + Marshal.PtrToStringAnsi(output));
        }
    }
}
