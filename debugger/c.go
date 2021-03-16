package debugger

// #cgo CFLAGS: -g -Wall
// #include <stdio.h>
// #include <stdlib.h>
// #include <windows.h>
//
///*	BOOL CreateProcessA(
//		LPCSTR                lpApplicationName,	//The name of the module to be executed
//		LPSTR                 lpCommandLine,		//The command line to be executed
//		LPSECURITY_ATTRIBUTES lpProcessAttributes,	//A pointer to a process SECURITY_ATTRIBUTES
//		LPSECURITY_ATTRIBUTES lpThreadAttributes,	//A pointer to a thread SECURITY_ATTRIBUTES
//		BOOL                  bInheritHandles,		//Inheritance flag
//		DWORD                 dwCreationFlags,		//The flags that control the priority class and the creation of the process
//		LPVOID                lpEnvironment,		//A pointer to the environment block for the new process
//		LPCSTR                lpCurrentDirectory,	//The full path to the current directory for the process
//		LPSTARTUPINFOA        lpStartupInfo,		//A pointer to a STARTUPINFO structure
//		LPPROCESS_INFORMATION lpProcessInformation	//A pointer to a PROCESS_INFORMATION structure
//	);*/
//
//struct RunProcRet {
//	PROCESS_INFORMATION	info;
//	DWORD error;
//};
//
//struct RunProcRet runProc(char *fname) {
//	PROCESS_INFORMATION  ProcessInformation;
//	STARTUPINFOA         StartupInfo;
//	struct RunProcRet ret = {ProcessInformation, 0};
// 	memset(&StartupInfo, 0, sizeof(StartupInfo));
//	StartupInfo.cb = sizeof(StartupInfo);
//	if (!CreateProcessA(
//  	(LPCSTR)fname,
// 		NULL, NULL, NULL, FALSE,
//		DEBUG_PROCESS, NULL, NULL,
//		(LPSTARTUPINFOA)&StartupInfo,
//		&ret.info)) {
//		ret.error =  GetLastError();
//	}
//	return ret;
// }
//
//
//
//struct ThreadContextRet {
//	CONTEXT context;
//	DWORD error;
//};
//
//struct ThreadContextRet threadContext(PROCESS_INFORMATION ProcessInformation) {
//	CONTEXT				 	tContext;
//	tContext.ContextFlags = CONTEXT_FULL;
//	struct ThreadContextRet ret = {tContext, 0};
//	if (!GetThreadContext(ProcessInformation.hThread, &tContext)) {
//		ret.error =  GetLastError();
//	}
//	return ret;
//}
//
//
//
//struct DebugEventRet {
//	DEBUG_EVENT event;
//	DWORD error;
//};
//
//struct DebugEventRet waitForDebugEvent() {
//	DEBUG_EVENT			 DebugEvent;
// 	struct DebugEventRet ret = {DebugEvent, 0};
//	if (!WaitForDebugEvent(&ret.event, 100)) {
//		if (GetLastError() != ERROR_SEM_TIMEOUT) {
//			ret.error =  GetLastError();
//		}
//	}
//	return ret;
//}
//
//
//
// struct ContinueDebugEventRet {
//	DEBUG_EVENT event;
//	int cont;
//};
//
//struct ContinueDebugEventRet continueDebugEvent(DEBUG_EVENT DebugEvent) {
//	struct ContinueDebugEventRet ret = {DebugEvent, 0};
//	if (DebugEvent.dwDebugEventCode != EXCEPTION_DEBUG_EVENT) {
//		ContinueDebugEvent(DebugEvent.dwProcessId, DebugEvent.dwThreadId, DBG_CONTINUE);
//		ret.cont = 1;
//	}
//	return ret;
//}
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type debugger struct {
	name       string
	debugEvent C.struct__DEBUG_EVENT
	process    C.struct__PROCESS_INFORMATION
	context    interface{}
}

func New(name string) Debugger {
	return &debugger{name: name}
}

func (d *debugger) processRun(name string) error {
	cs := C.CString(name)
	process := C.runProc(cs)
	if process.error != 0 {
		return errors.New(fmt.Sprintf("waitForDebugEvent exit code: %d", process.error))
	}
	C.free(unsafe.Pointer(cs))

	d.process = process.info
	return nil
}

func (d *debugger) waitForDebugEvent() error {
	debugEvent := C.waitForDebugEvent()
	if debugEvent.error != 0 {
		return errors.New(fmt.Sprintf("waitForDebugEvent exit code: %d", debugEvent.error))
	}
	d.debugEvent = debugEvent.event
	return nil
}

func (d *debugger) threadContext() error {
	threadContext := C.threadContext(d.process)
	if threadContext.error != 0 {
		return errors.New(fmt.Sprintf("threadContext exit code: %d", threadContext.error))
	}
	d.context = threadContext.context
	return nil
}

func (d *debugger) continueDebugEvent() bool {
	debugEvent := C.continueDebugEvent(d.debugEvent)
	d.debugEvent = debugEvent.event
	if debugEvent.cont == 0 {
		return false
	}
	return true
}
