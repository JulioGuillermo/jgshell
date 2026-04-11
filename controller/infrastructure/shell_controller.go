package controllerinfrastructure

import (
	"sync"

	autocompletedomain "github.com/julioguillermo/jgshell/autocomplete/domain"
	autocompleteinfrastructure "github.com/julioguillermo/jgshell/autocomplete/infrastructure"
	executorapplication "github.com/julioguillermo/jgshell/executor/application"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	persistencedomain "github.com/julioguillermo/jgshell/persistence/domain"
	persistenceinfrastructure "github.com/julioguillermo/jgshell/persistence/infrastructure"
	routerapplication "github.com/julioguillermo/jgshell/router/application"
	routerdomain "github.com/julioguillermo/jgshell/router/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	shellinfrastructure "github.com/julioguillermo/jgshell/shell/infrastructure"
	shelldetectorapplication "github.com/julioguillermo/jgshell/shelldetector/application"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
	statusdomain "github.com/julioguillermo/jgshell/status/domain"
	statusinfrastructure "github.com/julioguillermo/jgshell/status/infrastructure"
	toolsapplication "github.com/julioguillermo/jgshell/tools/application"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
	toolsinfrastructure "github.com/julioguillermo/jgshell/tools/infrastructure"
	wrapperapplication "github.com/julioguillermo/jgshell/wrapper/application"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
	wrapperinfrastructure "github.com/julioguillermo/jgshell/wrapper/infrastructure"
)

type ShellController struct {
	Locker sync.Locker

	UUIDGenerator  toolsdomain.UUIDGenerator
	OutputCleanner toolsdomain.OutputCleaner

	Persistencer persistencedomain.PersistenceCtl
	History      executordomain.History

	Router routerdomain.Router

	Shell               shelldomain.FullShell
	ShellDetector       shelldetectordomain.ShellDetector
	ShellWrapper        wrapperdomain.ShellWrapper
	ShellCmdWrapper     wrapperdomain.CmdWrapper
	ShellSimpleExecutor executordomain.SimpleExecutor
	ShellFastExecutor   executordomain.FastExecutor
	ShellExecutor       executordomain.Executor

	StatusLoader statusdomain.StatusLoader
	Autocomplete autocompletedomain.Autocomplete
}

func NewShellController(cmd string) (*ShellController, error) {
	persistenceCtl, err := persistenceinfrastructure.NewPersistenceCtl()
	if err != nil {
		return nil, err
	}

	ctl := &ShellController{
		Locker:         &sync.Mutex{},
		UUIDGenerator:  toolsinfrastructure.NewUUIDGenerator(),
		OutputCleanner: toolsapplication.NewOutputCleaner(),
		History:        executorapplication.NewHistory(),
		Persistencer:   persistenceCtl,
	}

	err = ctl.initShell(cmd)
	if err != nil {
		return nil, err
	}

	err = ctl.initExecutors()
	if err != nil {
		return nil, err
	}

	err = ctl.initFeatures()
	if err != nil {
		return nil, err
	}

	return ctl, nil
}

func (ctl *ShellController) initShell(cmd string) error {
	sh, err := shellinfrastructure.NewShellConnector(cmd)
	if err != nil {
		return err
	}
	ctl.Shell = sh
	router, err := routerapplication.NewRouter(
		ctl.Shell,
		routerapplication.NewQueue(
			executordomain.SimpleQueue,
			wrapperdomain.REWrapStartSimple,
			wrapperdomain.REWrapDoneSimple,
		),
		routerapplication.NewQueue(
			executordomain.FastQueue,
			wrapperdomain.REWrapStartFast,
			wrapperdomain.REWrapDone,
		),
		routerapplication.NewQueue(
			executordomain.DefaultQueue,
			wrapperdomain.REWrapStart,
			wrapperdomain.REWrapDone,
		),
	)
	if err != nil {
		return err
	}
	ctl.Router = router
	return nil
}

func (ctl *ShellController) initExecutors() error {
	ctl.ShellSimpleExecutor = executorapplication.NewSimpleExecutor(ctl.Router, ctl.UUIDGenerator)
	ctl.ShellDetector = shelldetectorapplication.NewShellDetector(ctl.ShellSimpleExecutor)
	ctl.ShellWrapper = wrapperinfrastructure.NewShellWrapper(ctl.Shell, ctl.ShellDetector)
	ctl.ShellCmdWrapper = wrapperapplication.NewCmdWrapper()
	ctl.ShellFastExecutor = executorapplication.NewFastExecutor(ctl.Router, ctl.ShellDetector, ctl.ShellCmdWrapper)
	ctl.ShellExecutor = executorapplication.NewExecutor(ctl.Router, ctl.ShellDetector, ctl.ShellCmdWrapper, ctl.UUIDGenerator)

	return nil
}

func (ctl *ShellController) initFeatures() error {
	ctl.StatusLoader = statusinfrastructure.NewStatusLoader(ctl.ShellDetector, ctl.ShellFastExecutor)
	ctl.Autocomplete = autocompleteinfrastructure.NewAutocomplete(ctl.ShellDetector, ctl.ShellFastExecutor)
	return nil
}
