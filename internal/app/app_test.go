package app

import (
	"fmt"
	"sync"
	"testing"

	"github.com/protosio/protos/internal/core"
	"github.com/protosio/protos/internal/mock"
	"github.com/protosio/protos/internal/util"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestAppManager(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rmMock := mock.NewMockResourceManager(ctrl)
	tmMock := mock.NewMockTaskManager(ctrl)
	rpMock := mock.NewMockRuntimePlatform(ctrl)
	dbMock := mock.NewMockDB(ctrl)
	wspMock := mock.NewMockWSPublisher(ctrl)
	cmMock := mock.NewMockCapabilityManager(ctrl)
	metaMock := mock.NewMockMeta(ctrl)
	pruMock := mock.NewMockPlatformRuntimeUnit(ctrl)
	appStoreMock := NewMockappStore(ctrl)
	capMock := mock.NewMockCapability(ctrl)

	c := make(chan interface{}, 10)

	// test app manager creation and initial app loading from db
	dbMock.EXPECT().All(gomock.Any()).Return(nil).Times(1).
		Do(func(to interface{}) {
			apps := to.(*[]*App)
			*apps = append(*apps,
				&App{ID: "id1", Name: "app1", Status: statusUnknown, Tasks: []string{"1", "2"}, access: &sync.Mutex{}, PublicPorts: []util.Port{util.Port{Nr: 10000, Type: util.TCP}}},
				&App{ID: "id2", Name: "app2", Status: statusUnknown, Tasks: []string{"1"}, access: &sync.Mutex{}},
				&App{ID: "id3", Name: "app3", Status: statusUnknown, Tasks: []string{"1"}, access: &sync.Mutex{}})
		})

	// one of the inputs is nil
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("A nil input in the CreateManager call should lead to a panic")
			}
		}()
		CreateManager(rmMock, nil, rpMock, dbMock, metaMock, wspMock, nil, cmMock)
	}()

	// happy case
	am := CreateManager(rmMock, tmMock, rpMock, dbMock, metaMock, wspMock, appStoreMock, cmMock)

	//
	// GetCopy
	//

	t.Run("GetCopy", func(t *testing.T) {
		_, err := am.GetCopy("wrongId")
		if err == nil {
			t.Errorf("GetCopy(wrongId) should return an error")
		}

		app, err := am.GetCopy("id1")
		if err != nil {
			t.Errorf("GetCopy(id1) should NOT return an error: %s", err.Error())
		} else {
			if app.GetName() != "app1" {
				t.Errorf("App id 'id1' should have name app1, NOT %s", app.GetName())
			}
		}
	})

	//
	// CopyAll
	//

	t.Run("CopyAll", func(t *testing.T) {
		if len(am.CopyAll()) != 3 {
			t.Errorf("CopyAll should return 3 apps. Instead it returned %d", len(am.CopyAll()))
		}
	})

	//
	// Read
	//

	t.Run("Read", func(t *testing.T) {
		_, err := am.Read("wrongId")
		if err == nil {
			t.Errorf("Read(wrongId) should return an error")
		}

		app, err := am.Read("id1")
		if err != nil {
			t.Errorf("Read(id1) should NOT return an error: %s", err.Error())
		} else {
			if app.GetName() != "app1" {
				t.Errorf("App id 'id1' should have name app1, NOT %s", app.GetName())
			}
		}
	})

	//
	// Select
	//
	t.Run("Select", func(t *testing.T) {
		filter := func(app core.App) bool {
			if app.GetName() == "app2" {
				return true
			}
			return false
		}

		apps := am.Select(filter)
		if len(apps) != 1 {
			t.Errorf("Select(filter) should return 1 app. Instead it returned %d", len(apps))
		}
		for _, app := range apps {
			if app.GetName() != "app2" || app.GetID() != "id2" {
				t.Errorf("Expected app id '%s' and app name '%s', but found '%s' and '%s'", "id2", "app2", app.GetID(), app.GetName())
			}
		}
	})

	//
	// CreateAsync
	//

	t.Run("CreateAsync", func(t *testing.T) {
		// one of the required inputs is nil
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("An empty required input in the CreateAsync call should lead to a panic")
				}
			}()
			am.CreateAsync("", "0.0.1", "c", &core.InstallerMetadata{}, map[string]string{}, false)
		}()
		tmMock.EXPECT().New(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		_ = am.CreateAsync("a", "b", "c", &core.InstallerMetadata{}, map[string]string{}, false)
	})

	//
	// Create
	//

	t.Run("Create", func(t *testing.T) {
		_, err := am.Create("a", "b", "", map[string]string{}, core.InstallerMetadata{}, "taskid")
		if err == nil {
			t.Errorf("Creating an app using a blank name should result in an error")
		}

		// installer params test
		_, err = am.Create("a", "b", "c", map[string]string{}, core.InstallerMetadata{Params: []string{"test"}}, "taskid")
		if err == nil {
			t.Errorf("Creating an app and not providing the mandatory params should result in an error")
		}

		// capability test, error while creating DNS for app
		metaMock.EXPECT().GetPublicIP().Return("1.1.1.1").Times(1)
		rmMock.EXPECT().CreateDNS(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test error")).Times(1)
		cmMock.EXPECT().GetByName("PublicDNS").Return(capMock, nil).Times(2)
		capMock.EXPECT().GetName().Return("PublicDNS").Times(1)
		cmMock.EXPECT().Validate(capMock, gomock.Any()).Return(true).Times(1)
		_, err = am.Create("a", "b", "c", map[string]string{}, core.InstallerMetadata{Capabilities: []string{"PublicDNS"}}, "taskid")
		if err == nil {
			t.Error("Creating an app and having a DNS creation error should result in an error")
		}

		// happy case
		wspMock.EXPECT().GetWSPublishChannel().Return(c).Times(1)
		tmMock.EXPECT().GetIDs(gomock.Any()).Return(*linkedhashmap.New()).Times(1)
		rmMock.EXPECT().Select(gomock.Any()).Return(map[string]core.Resource{}).Times(1)
		cmMock.EXPECT().GetByName("PublicDNS").Return(capMock, nil).Times(1)
		capMock.EXPECT().GetName().Return("PublicDNS").Times(1)
		dbMock.EXPECT().Save(gomock.Any()).Return(nil).Times(1)
		app, err := am.Create("a", "b", "c", map[string]string{}, core.InstallerMetadata{}, "taskid")

		dbMock.EXPECT().Remove(gomock.Any()).Return(nil).Times(1)
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(1)
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		am.Remove(app.ID)
	})

	//
	// GetServices
	//
	t.Run("GetServices", func(t *testing.T) {
		dnsRscType := NewMockdnsResource(ctrl)
		dnsRscType.EXPECT().GetName().Return("app1").Times(2)
		dnsRscType.EXPECT().GetValue().Return("1.1.1.1").Times(1)
		rscMock := mock.NewMockResource(ctrl)
		rscMock.EXPECT().GetAppID().Return("id1").Times(1)
		rscMock.EXPECT().GetValue().Return(dnsRscType).Times(1)
		rmMock.EXPECT().Select(gomock.Any()).Return(map[string]core.Resource{"1": rscMock}).Times(1)
		metaMock.EXPECT().GetDomain().Return("giurgiu.io").Times(1)

		services := am.GetServices()
		if len(services) != 1 {
			t.Fatalf("GetServices should only return 1 service in this test, but %d were returned", len(services))
		}
		svc := services[0]
		if len(svc.Ports) != 1 ||
			svc.Ports[0].Nr != 10000 ||
			svc.Ports[0].Type != util.TCP ||
			svc.Name != "app1" ||
			svc.Domain != "app1.giurgiu.io" {
			t.Error("GetServices returned a service with incorrect values")
		}
	})

	//
	// Remove
	//
	t.Run("Remove", func(t *testing.T) {
		initialNr := len(am.apps.apps)
		// non-existent app id
		err := am.Remove("id4")
		if err == nil {
			t.Error("Remove(id4) should return an error because id4 does not exist")
		}
		if len(am.apps.apps) != initialNr {
			t.Errorf("Wrong number of apps found: %d vs %d", len(am.apps.apps), initialNr)
		}

		// existent app id which returns an error on app.remove()
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(nil, fmt.Errorf("test error")).Times(1)
		dbMock.EXPECT().Remove(gomock.Any()).Return(nil).Times(1)
		err = am.Remove("id3")
		if err == nil {
			t.Error("Remove(id3) should return an error because app.remove() returns an error")
		}
		if len(am.apps.apps) != initialNr {
			t.Errorf("Wrong number of apps found: %d vs %d", len(am.apps.apps), initialNr)
		}

		// existent app id - happy path
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(1)
		dbMock.EXPECT().Remove(gomock.Any()).Return(nil).Times(1)
		err = am.Remove("id2")
		if err != nil {
			t.Errorf("Remove(id2) should NOT return an error: %s", err.Error())
		}
		if len(am.apps.apps) != initialNr-1 {
			t.Errorf("Wrong number of apps found: %d vs %d", len(am.apps.apps), initialNr-1)
		}
	})

	//
	// RemoveAsync
	//
	t.Run("RemoveAsync", func(t *testing.T) {
		taskMock := mock.NewMockTask(ctrl)
		tmMock.EXPECT().New(gomock.Any(), gomock.Any()).Return(taskMock).Times(1)
		removeTask := am.RemoveAsync("id4")
		if removeTask != taskMock {
			t.Error("RemoveAsync returned an incorrect task")
		}
	})

	//
	// saveApp
	//
	t.Run("saveApp", func(t *testing.T) {
		app2 := &App{ID: "id2", Name: "app2", access: &sync.Mutex{}, parent: am}
		wspMock.EXPECT().GetWSPublishChannel().Return(c).Times(2)
		pruMock.EXPECT().GetStatus().Return("exited").Times(2)
		pruMock.EXPECT().GetExitCode().Return(0).Times(2)
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(2)
		tmMock.EXPECT().GetIDs(gomock.Any()).Return(*linkedhashmap.New()).Times(2)
		rmMock.EXPECT().Select(gomock.Any()).Return(map[string]core.Resource{}).Times(2)

		// happy path
		dbMock.EXPECT().Save(gomock.Any()).Return(nil).Times(1)
		am.saveApp(app2)

		// db error should lead to panic
		dbMock.EXPECT().Save(gomock.Any()).Return(errors.New("test db error")).Times(1)
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("A DB error in saveApp should lead to a panic")
				}
			}()
			am.saveApp(app2)
		}()
	})

	//
	// CreateDevApp
	//
	t.Run("CreateDevApp", func(t *testing.T) {
		// app creation returns error
		_, err := am.CreateDevApp("a", "b", "", core.InstallerMetadata{}, map[string]string{})
		if err == nil {
			t.Error("CreateDevApp should fail when the app creation step fails")
		}

		// happy case
		wspMock.EXPECT().GetWSPublishChannel().Return(c).Times(2)
		tmMock.EXPECT().GetIDs(gomock.Any()).Return(*linkedhashmap.New()).Times(2)
		rmMock.EXPECT().Select(gomock.Any()).Return(map[string]core.Resource{}).Times(2)
		cmMock.EXPECT().GetByName("PublicDNS").Return(capMock, nil).Times(1)
		capMock.EXPECT().GetName().Return("PublicDNS").Times(1)
		dbMock.EXPECT().Save(gomock.Any()).Return(nil).Times(2)
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(1)
		pruMock.EXPECT().GetStatus().Return("exited").Times(1)
		pruMock.EXPECT().GetExitCode().Return(0).Times(1)
		app, err := am.CreateDevApp("a", "b", "c", core.InstallerMetadata{}, map[string]string{})
		if err != nil {
			t.Errorf("CreateDevApp(...) should NOT return an error: %s", err.Error())
		}
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(1)
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		am.Remove(app.GetID())
	})

	//
	// GetAllPublic
	//

	t.Run("GetAllPublic", func(t *testing.T) {
		nrOfApps := len(am.apps.apps)
		tasks := linkedhashmap.New()
		tasks.Put("1", gomock.Any())
		tasks.Put("2", gomock.Any())
		tmMock.EXPECT().GetAll().Return(tasks).Times(1)
		rpMock.EXPECT().GetDockerContainer(gomock.Any()).Return(pruMock, nil).Times(nrOfApps)
		pruMock.EXPECT().GetStatus().Return("exited").Times(nrOfApps)
		pruMock.EXPECT().GetExitCode().Return(0).Times(nrOfApps)
		papps := am.GetAllPublic()
		if len(papps) != nrOfApps {
			t.Errorf("GetAllPublic() should return %d apps, but it returned %d", nrOfApps, len(papps))
		}
		tsks1 := linkedhashmap.Map(papps["id1"].(*PublicApp).Tasks)
		if len(tsks1.Keys()) != 2 {
			t.Errorf("There should be 2 tasks in the public app with id1, but there are %d", len(tsks1.Keys()))
		}
		tsks2 := linkedhashmap.Map(papps["id3"].(*PublicApp).Tasks)
		if len(tsks2.Keys()) != 1 {
			t.Errorf("There should be 1 tasks in the public app with id2, but there are %d", len(tsks2.Keys()))
		}

	})

}

func TestApp(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parentMock := NewMockappParent(ctrl)
	platformMock := mock.NewMockRuntimePlatform(ctrl)
	tmMock := mock.NewMockTaskManager(ctrl)
	pruMock := mock.NewMockPlatformRuntimeUnit(ctrl)
	taskMock := mock.NewMockTask(ctrl)
	rmMock := mock.NewMockResourceManager(ctrl)
	capMock := mock.NewMockCapability(ctrl)
	cmMock := mock.NewMockCapabilityManager(ctrl)

	app := &App{
		ID:          "id1",
		Name:        "app1",
		Status:      "initial",
		parent:      parentMock,
		access:      &sync.Mutex{},
		PublicPorts: []util.Port{util.Port{Nr: 10000, Type: util.TCP}}}

	//
	// GetID
	//

	t.Run("GetID", func(t *testing.T) {
		if app.GetID() != "id1" {
			t.Error("GetID should return id1")
		}
	})

	//
	// GetName
	//

	t.Run("GetName", func(t *testing.T) {
		if app.GetName() != "app1" {
			t.Error("GetName should return app1")
		}
	})

	//
	// SetStatus
	//

	t.Run("SetStatus", func(t *testing.T) {
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		teststatus := "teststatus"
		app.SetStatus(teststatus)
		if app.Status != teststatus {
			t.Errorf("SetStatus did not set the correct status. Status should be %s but is %s", teststatus, app.Status)
		}
	})

	//
	// AddAction
	//

	t.Run("AddAction", func(t *testing.T) {
		_, err := app.AddAction("bogus")
		if err == nil {
			t.Error("AddAction(bogus) should fail and return an error")
		}

		parentMock.EXPECT().getTaskManager().Return(tmMock).Times(2)
		tmMock.EXPECT().New(gomock.Any(), gomock.Any()).Return(taskMock).Times(2)
		tsk, err := app.AddAction("start")
		if err != nil {
			t.Error("AddAction(start) should NOT return an error")
		}
		if tsk != taskMock {
			t.Error("AddAction(start) returned an incorrect task")
		}
		tsk, err = app.AddAction("stop")
		if err != nil {
			t.Error("AddAction(stop) should NOT return an error")
		}
		if tsk != taskMock {
			t.Error("AddAction(stop) returned an incorrect task")
		}
	})

	//
	// AddTask
	//

	t.Run("AddTask", func(t *testing.T) {
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		app.AddTask("tskid")
		if present, _ := util.StringInSlice("tskid", app.Tasks); present == false {
			t.Error("AddTask(tskid) did not lead to 'tskid' being present in the Tasks slice")
		}
	})

	//
	// Save
	//

	t.Run("Save", func(t *testing.T) {
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		app.Save()
	})

	//
	// createContainer
	//

	t.Run("createContainer", func(t *testing.T) {
		app.InstallerMetadata.PersistancePath = ""
		app.InstallerMetadata.PersistancePath = "/data"
		// volume creation error
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetOrCreateVolume(gomock.Any(), gomock.Any()).Return("volumeid", errors.New("volume error")).Times(1)
		_, err := app.createContainer()
		if err == nil {
			t.Error("createContainer should return an error when the volume creation errors out")
		}

		// new container error
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		platformMock.EXPECT().GetOrCreateVolume(gomock.Any(), gomock.Any()).Return("volumeid", nil).Times(1)
		platformMock.EXPECT().NewContainer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("container error")).Times(1)
		_, err = app.createContainer()
		if err == nil {
			t.Error("createContainer should return an error when the container creation errors out")
		}

		// happy case
		pruMock.EXPECT().GetID().Return("cntid").Times(1)
		pruMock.EXPECT().GetIP().Return("cntip").Times(1)
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetOrCreateVolume(gomock.Any(), gomock.Any()).Return("volumeid", nil).Times(1)
		platformMock.EXPECT().NewContainer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pruMock, nil).Times(1)
		_, err = app.createContainer()
		if err != nil {
			t.Errorf("createContainer should NOT return an error: %s", err.Error())
		}
	})

	//
	// getOrCreateContainer
	//

	t.Run("getOrCreateContainer", func(t *testing.T) {
		app.InstallerMetadata.PersistancePath = ""

		// container retrieval error
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, errors.New("container retrieval error"))
		_, err := app.getOrCreateContainer()
		if err == nil {
			t.Error("getOrCreateContainer() should return an error when the container can't be retrieved")
		}

		// container retrieval returns err of type core.ErrContainerNotFound, and container creation fails
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, util.NewTypedError("container retrieval error", core.ErrContainerNotFound))
		platformMock.EXPECT().NewContainer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("container creation error")).Times(1)
		_, err = app.getOrCreateContainer()
		if err == nil {
			t.Error("getOrCreateContainer() should return an error when no container exists and the creation of one fails")
		}

		// container retrieval returns err and creation of a new container works
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, util.NewTypedError("container retrieval error", core.ErrContainerNotFound))
		platformMock.EXPECT().NewContainer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pruMock, nil).Times(1)
		pruMock.EXPECT().GetID().Return("cntid").Times(1)
		pruMock.EXPECT().GetIP().Return("cntip").Times(1)
		cnt, err := app.getOrCreateContainer()
		if err != nil {
			t.Errorf("getOrCreateContainer() should not return an error: %s", err.Error())
		}
		if cnt != pruMock {
			t.Errorf("getOrCreateContainer() returned an incorrect container: %p vs %p", cnt, pruMock)
		}

		// container retrieval works
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		cnt, err = app.getOrCreateContainer()
		if err != nil {
			t.Errorf("getOrCreateContainer() should not return an error: %s", err.Error())
		}
		if cnt != pruMock {
			t.Errorf("getOrCreateContainer() returned an incorrect container: %p' vs %p", cnt, pruMock)
		}
	})

	//
	// enrichAppData
	//

	t.Run("enrichAppData", func(t *testing.T) {
		// app is creating, nothing is done
		app.Status = statusCreating
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(0)
		app.enrichAppData()
		if app.Status != statusCreating {
			t.Errorf("enrichAppData failed. App status should be '%s' but is '%s'", statusCreating, app.Status)
		}

		// app failes to retrieve container
		app.Status = "test"
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, errors.New("container retrieval error"))
		app.enrichAppData()
		if app.Status != statusUnknown {
			t.Errorf("enrichAppData failed. App status should be '%s' but is '%s'", statusUnknown, app.Status)
		}

		// app failes to retrieve container because the container is not found
		app.Status = "test"
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, util.NewTypedError("container retrieval error", core.ErrContainerNotFound))
		app.enrichAppData()
		if app.Status != statusStopped {
			t.Errorf("enrichAppData failed. App status should be '%s' but is '%s'", statusStopped, app.Status)
		}

		// app retrieves container and status is updates based on the container
		app.Status = "test"
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().GetStatus().Return("exited").Times(1)
		pruMock.EXPECT().GetExitCode().Return(1).Times(1)
		app.enrichAppData()
		if app.Status != statusFailed {
			t.Errorf("enrichAppData failed. App status should be '%s' but is '%s'", statusStopped, app.Status)
		}
	})

	//
	// StartAsync
	//

	t.Run("StartAsync", func(t *testing.T) {
		parentMock.EXPECT().getTaskManager().Return(tmMock).Times(1)
		tmMock.EXPECT().New(gomock.Any(), gomock.Any()).Return(taskMock).Times(1)
		tsk := app.StartAsync()
		if tsk != taskMock {
			t.Errorf("StartAsync() returned an incorrect task: %p vs %p", tsk, taskMock)
		}
	})

	//
	// Start
	//

	t.Run("Start", func(t *testing.T) {
		// failed container retrieval
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, errors.New("container retrieval error"))
		err := app.Start()
		if err == nil {
			t.Error("Start() should return an error when the container can't be retrieved")
		}
		if app.Status != statusFailed {
			t.Errorf("App status on failed start should be '%s' but is '%s'", statusFailed, app.Status)
		}

		// container failes to start
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Start().Return(errors.New("failed to start test container"))
		err = app.Start()
		if err == nil {
			t.Error("Start() should return an error when the container can't be started")
		}
		if app.Status != statusFailed {
			t.Errorf("App status on failed start should be '%s' but is '%s'", statusFailed, app.Status)
		}

		// happy case
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Start().Return(nil)
		err = app.Start()
		if err != nil {
			t.Errorf("Start() should not return an error: %s", err.Error())
		}
		if app.Status != statusRunning {
			t.Errorf("App status on successful start should be '%s' but is '%s'", statusRunning, app.Status)
		}
	})

	//
	// StopAsync
	//

	t.Run("StopAsync", func(t *testing.T) {
		parentMock.EXPECT().getTaskManager().Return(tmMock).Times(1)
		tmMock.EXPECT().New(gomock.Any(), gomock.Any()).Return(taskMock).Times(1)
		tsk := app.StopAsync()
		if tsk != taskMock {
			t.Errorf("StopAsync() returned an incorrect task: %p vs %p", tsk, taskMock)
		}
	})

	//
	// Stop
	//

	t.Run("Stop", func(t *testing.T) {
		// failed container retrieval
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, errors.New("container retrieval error"))
		err := app.Stop()
		if err == nil {
			t.Error("Stop() should return an error when the container can't be retrieved")
		}
		if app.Status != statusUnknown {
			t.Errorf("App status when Stop() fails (because of failure to retrieve container) should be '%s' but is '%s'", statusUnknown, app.Status)
		}

		// container not found
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, util.NewTypedError("container retrieval error", core.ErrContainerNotFound))
		err = app.Stop()
		if err != nil {
			t.Error("Stop() should NOT return an error when the container of an app does not exist")
		}
		if app.Status != statusStopped {
			t.Errorf("App status when Stop() succeeds (and app container does not exist) should be '%s' but is '%s'", statusStopped, app.Status)
		}

		// container fails to stop
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Stop().Return(errors.New("failed to stop container"))
		err = app.Stop()
		if err == nil {
			t.Error("Stop() should return an error when the container can't be stopped")
		}
		if app.Status != statusUnknown {
			t.Errorf("App status on unsuccessful stop should be '%s' but is '%s'", statusUnknown, app.Status)
		}

		// happy case
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Return().Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Stop().Return(nil)
		err = app.Stop()
		if err != nil {
			t.Errorf("Stop() should NOT return an error: %s", err.Error())
		}
		if app.Status != statusStopped {
			t.Errorf("App status on successful stop should be '%s' but is '%s'", statusStopped, app.Status)
		}
	})

	//
	// remove
	//

	t.Run("remove", func(t *testing.T) {
		// can't retrieve container
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, errors.New("container retrieval error"))
		err := app.remove()
		if err == nil {
			t.Error("remove() should return an error when the container can't be retrieved")
		}

		// container not found
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(nil, util.NewTypedError("container retrieval error", core.ErrContainerNotFound))
		err = app.remove()
		if err != nil {
			t.Errorf("remove() should NOT return an error when the container is not found: %s", err.Error())
		}

		// container retrieved and failed to remove it
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Remove().Return(errors.New("container removal error")).Times(1)
		err = app.remove()
		if err == nil {
			t.Error("remove() should return an error when the container can't be removed")
		}

		// container retrieved and removed
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		err = app.remove()
		if err != nil {
			t.Errorf("remove() should NOT return an error when the container is removed successfully: %s", err.Error())
		}

		// failed to remove volume
		app.VolumeID = "testvol"
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		platformMock.EXPECT().RemoveVolume(app.VolumeID).Return(errors.New("volume removal error"))
		err = app.remove()
		if err == nil {
			t.Error("remove() should return an error when the volume can't be removed")
		}

		// volume removed
		app.VolumeID = "testvol"
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(2)
		platformMock.EXPECT().GetDockerContainer("cntid").Return(pruMock, nil)
		pruMock.EXPECT().Remove().Return(nil).Times(1)
		platformMock.EXPECT().RemoveVolume(app.VolumeID).Return(nil)
		err = app.remove()
		if err != nil {
			t.Errorf("remove() should NOT return an error when the is removed successfully: %s", err.Error())
		}
	})

	//
	// ReplaceContainer
	//

	t.Run("ReplaceContainer", func(t *testing.T) {
		// container can't be retrieved
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("newcntid").Return(nil, errors.New("container retrieval error"))
		err := app.ReplaceContainer("newcntid")
		if err == nil {
			t.Error("ReplaceContainer() should return an error when the container can't be retrieved")
		}

		// happy case
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer("newcntid").Return(pruMock, nil)
		pruMock.EXPECT().GetIP().Return("1.1.1.1").Times(1)
		parentMock.EXPECT().saveApp(gomock.Any()).Times(1)
		err = app.ReplaceContainer("newcntid")
		if err != nil {
			t.Errorf("ReplaceContainer() should NOT return an error: %s", err.Error())
		}
		if app.ContainerID != "newcntid" || app.GetIP() != "1.1.1.1" {
			t.Errorf("ReplaceContainer() should lead to different data in app struct: %#v", app)
		}
	})

	//
	// GetIP
	//

	t.Run("GetIP", func(t *testing.T) {
		app.IP = "1.1.1.1"
		ip := app.GetIP()
		if ip != app.IP {
			t.Errorf("GetIP() returned an incorrect IP address. Should be '%s' but is '%s'", app.IP, ip)
		}
	})

	//
	// message queue related tests
	//

	t.Run("Message queue", func(t *testing.T) {
		//
		// SetMsgQ
		//
		wsc := &core.WSConnection{Close: make(chan bool, 1), Send: make(chan interface{}, 1)}
		app.SetMsgQ(wsc)
		if app.msgq != wsc {
			t.Error("SetMsgQ() set an incorrect msgq on the app struct")
		}

		//
		// CloseMsgQ
		//

		app.CloseMsgQ()
		if <-wsc.Close != true {
			t.Error("CloseMsgQ() did not lead to a close message in the Close channel")
		}

		//
		// SendMsg
		//

		// msgq is nil
		err := app.SendMsg("test")
		if err == nil {
			t.Error("SendMsg() should return an error when the app msgq is nil")
		}

		// happy case
		app.SetMsgQ(wsc)
		msg1 := "qmsg"
		err = app.SendMsg(msg1)
		if err != nil {
			t.Errorf("SendMsg() should not return an error when the msgq is set: %s", err.Error())
		}
		msg2 := <-wsc.Send
		if msg2 != msg2 {
			t.Errorf("SendMsg() sent an incorrect message: '%s' vs '%s'", msg2, msg1)
		}
	})

	//
	// CreateResource
	//

	t.Run("CreateResource", func(t *testing.T) {
		rscPayload := []byte("payload")
		rscMock := mock.NewMockResource(ctrl)

		// error creating rsc by ResourceManager
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().CreateFromJSON(rscPayload, app.ID).Return(nil, errors.New("failed to create resource"))
		_, err := app.CreateResource(rscPayload)
		if err == nil {
			t.Error("CreateResource() should return an error when the resource manager fails to create the resource")
		}

		// happy case
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().CreateFromJSON(rscPayload, app.ID).Return(rscMock, nil)
		rscMock.EXPECT().GetID().Return("rscid").Times(1)
		parentMock.EXPECT().saveApp(app).Return().Times(1)
		rsc, err := app.CreateResource(rscPayload)
		if err != nil {
			t.Errorf("CreateResource() should NOT return an error: %s", err.Error())
		}
		if rsc != rscMock {
			t.Errorf("CreateResource() returned an incorrect resource instance: %p vs %p", rsc, rscMock)
		}
	})

	//
	// DeleteResource
	//

	t.Run("DeleteResource", func(t *testing.T) {
		app.Resources = []string{"rscid"}

		// rsc id not owned by app
		err := app.DeleteResource("rscid2")
		if err == nil {
			t.Error("DeleteResource() should return an error when the resource id is not owned by the app")
		}

		// rsc can't be deleted by the resource manager
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Delete("rscid").Return(errors.New("can't delete resource"))
		err = app.DeleteResource("rscid")
		if err == nil {
			t.Error("DeleteResource() should return an error when the resource can't be deleted by the resource manager")
		}

		// happy case
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Delete("rscid").Return(nil)
		parentMock.EXPECT().saveApp(app).Return().Times(1)
		err = app.DeleteResource("rscid")
		if err != nil {
			t.Errorf("DeleteResource() should NOT return an error: %s", err.Error())
		}
		if len(app.Resources) != 0 {
			t.Error("DeleteResource() did not correctly remove the resource id from the app struct")
		}
	})

	//
	// GetResources
	//

	t.Run("GetResources", func(t *testing.T) {
		// 0 resources
		rscids := []string{}
		app.Resources = rscids
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rscs := app.GetResources()
		if len(rscs) != len(rscids) {
			t.Errorf("GetResources() should return %d resources but it returned %d", len(rscids), len(rscs))
		}

		// 2 resources, 1 is not retrieved by the resource manager
		rscids = []string{"rscid1", "rscid2"}
		rsc1 := mock.NewMockResource(ctrl)
		app.Resources = rscids
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Get(rscids[0]).Return(rsc1, nil)
		rmMock.EXPECT().Get(rscids[1]).Return(nil, errors.New("couldn't retrieve resource"))
		rscs = app.GetResources()
		if len(rscs) != 1 {
			t.Errorf("GetResources() should return %d resources but it returned %d", 1, len(rscs))
		}
		if rscs[rscids[0]] != rsc1 {
			t.Errorf("GetResources() returned the wrong resource. Id should be %s but is %s", rscids[0], rsc1.GetID())
		}

		// happy case, 2 resources, both successfully retrieved
		rscids = []string{"rscid1", "rscid2"}
		rsc1 = mock.NewMockResource(ctrl)
		rsc2 := mock.NewMockResource(ctrl)
		app.Resources = rscids
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Get(rscids[0]).Return(rsc1, nil)
		rmMock.EXPECT().Get(rscids[1]).Return(rsc2, nil)
		rscs = app.GetResources()
		if len(rscs) != len(rscids) {
			t.Errorf("GetResources() should return %d resources but it returned %d", len(rscids), len(rscs))
		}
		if rscs[rscids[0]] != rsc1 || rscs[rscids[1]] != rsc2 {
			t.Error("GetResources() returned the wrong resources")
		}
	})

	//
	// GetResource
	//

	t.Run("GetResource", func(t *testing.T) {
		// bad resource id which does not belong to the app
		_, err := app.GetResource("rscid3")
		if err == nil {
			t.Error("GetResource() should return an error when the resource id does not belong to the app")
		}

		// correct resource ID but resource manager can't retrieve it
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Get("rscid2").Return(nil, errors.New("couldn't retrieve resource"))
		_, err = app.GetResource("rscid2")
		if err == nil {
			t.Error("GetResource() should return an error when the resource id cannot be retrieved by the resource manager")
		}

		// happy case
		rsc1 := mock.NewMockResource(ctrl)
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Get("rscid2").Return(rsc1, nil)
		rsc, err := app.GetResource("rscid2")
		if err != nil {
			t.Errorf("GetResource() should NOT return an error: %s", err.Error())
		}
		if rsc != rsc1 {
			t.Errorf("GetResource() returned an incorrect resource: %p vs %p", rsc, rsc1)
		}
	})

	//
	// ValidateCapability
	//

	t.Run("ValidateCapability", func(t *testing.T) {
		// app doesn't have AuthUser capability
		app.Capabilities = []string{}
		capMock.EXPECT().GetName().Return("AuthUser").Times(1)
		err := app.ValidateCapability(capMock)
		if err == nil {
			t.Error("ValidateCapability() should return an error when the app doesn't have that capability")
		}

		// app has AuthUser capability
		app.Capabilities = []string{"AuthUser"}
		parentMock.EXPECT().getCapabilityManager().Return(cmMock).Times(1)
		cmMock.EXPECT().Validate(capMock, "AuthUser").Return(true).Times(1)
		err = app.ValidateCapability(capMock)
		if err != nil {
			t.Errorf("ValidateCapability() should NOT return an error when the app has that capability: %s", err.Error())
		}
	})
	//
	// Provides
	//

	t.Run("Provides", func(t *testing.T) {
		if app.Provides("non-existent rsc type") != false {
			t.Error("Provides() should return false when called with a rsc type that is not provided by the app")
		}
		app.InstallerMetadata.Provides = []string{"rsctype"}
		if app.Provides("rsctype") != true {
			t.Error("Provides() should return true when called with a rsc type that is provided by the app")
		}
	})

	//
	// Public
	//

	t.Run("Public", func(t *testing.T) {
		app.Tasks = []string{"1"}
		app.Resources = []string{"1"}
		tasks := linkedhashmap.New()
		tasks.Put("1", gomock.Any())
		resources := map[string]core.Resource{"1": mock.NewMockResource(ctrl)}
		parentMock.EXPECT().getPlatform().Return(platformMock).Times(1)
		platformMock.EXPECT().GetDockerContainer(app.ContainerID).Return(pruMock, nil).Times(1)
		pruMock.EXPECT().GetStatus().Return("yolo").Times(1)
		pruMock.EXPECT().GetExitCode().Return(0).Times(1)
		parentMock.EXPECT().getTaskManager().Return(tmMock).Times(1)
		tmMock.EXPECT().GetIDs(app.Tasks).Return(*tasks).Times(1)
		parentMock.EXPECT().getResourceManager().Return(rmMock).Times(1)
		rmMock.EXPECT().Select(gomock.Any()).Return(resources).Times(1)

		papp := app.Public()
		ptasks := linkedhashmap.Map(papp.(*PublicApp).Tasks)
		_, found := ptasks.Get("1")
		if found == false {
			t.Error("Public() returned a public app that does not contain the correct tasks")
		}
		presources := papp.(*PublicApp).Resources
		if presources["1"] != resources["1"] {
			t.Errorf("Public() returned a public app with an incorrect resources var: %p vs %p", presources["1"], resources["1"])
		}
	})

}

func TestTask(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amMock := NewMocktaskParent(ctrl)

	t.Run("CreateAppTask", func(t *testing.T) {
		task := CreateAppTask{
			am:                nil,
			InstallerID:       "1",
			InstallerVersion:  "",
			AppName:           "testapp",
			InstallerMetadata: nil,
			InstallerParams:   map[string]string{},
			StartOnCreation:   false,
		}
		tskID := "1"
		p := mock.NewMockProgress(ctrl)
		store := NewMockappStore(ctrl)
		inst := mock.NewMockInstaller(ctrl)
		app := NewMockapp(ctrl)
		baseTaskMock := mock.NewMockTask(ctrl)
		downloadTaskMock := mock.NewMockTask(ctrl)
		startAsyncTaskMock := mock.NewMockTask(ctrl)
		pruMock := mock.NewMockPlatformRuntimeUnit(ctrl)

		//
		// Run
		//

		// required inputs are missing for the task
		err := task.Run(baseTaskMock, tskID, p)
		log.Info(err)
		if err == nil {
			t.Error("Run() should return an error when one of the required task fields is empty")
		}
		task.am = amMock
		task.InstallerVersion = "0.0.0-dev"

		// failed to get installer metadata from the store
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(nil, errors.New("failed to retrieve image")).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the retrieval of the installer from the store fails")
		}

		// failed to retrieve installer metadata
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, errors.New("failed to retrieve install metadata")).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the retrieval of the installer metadata fails")
		}

		// app manager fails to create app
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, nil).Times(1)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(nil, errors.New("failed to create app")).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the app manager fails to create the app")
		}

		// image not available locally and download task returns an error
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, nil).Times(1)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(app, nil).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		p.EXPECT().SetPercentage(10).Times(1)
		p.EXPECT().SetState("Created application").Times(1)
		inst.EXPECT().IsPlatformImageAvailable(task.InstallerVersion).Return(false, nil).Times(1)
		app.EXPECT().GetID().Return("appid1").Times(1)
		inst.EXPECT().DownloadAsync(task.InstallerVersion, "appid1").Return(downloadTaskMock).Times(1)
		downloadTaskMock.EXPECT().GetID().Return("2")
		app.EXPECT().AddTask("2").Times(1)
		downloadTaskMock.EXPECT().Wait().Return(errors.New("download task error"))
		app.EXPECT().SetStatus(statusFailed).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the download image task fails")
		}

		// image available locally and create container fails
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, nil).Times(1)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(app, nil).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		p.EXPECT().SetPercentage(10).Times(1)
		p.EXPECT().SetState("Created application").Times(1)
		inst.EXPECT().IsPlatformImageAvailable(task.InstallerVersion).Return(true, nil).Times(1)
		p.EXPECT().SetPercentage(50).Times(1)
		p.EXPECT().SetState("Docker image found locally").Times(1)
		app.EXPECT().createContainer().Return(nil, errors.New("failed to create container"))
		app.EXPECT().SetStatus(statusFailed).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the app container fails to be created")
		}

		// start on creation is true and app fails to start
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, nil).Times(1)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(app, nil).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		p.EXPECT().SetPercentage(10).Times(1)
		p.EXPECT().SetState("Created application").Times(1)
		inst.EXPECT().IsPlatformImageAvailable(task.InstallerVersion).Return(true, nil).Times(1)
		p.EXPECT().SetPercentage(50).Times(1)
		p.EXPECT().SetState("Docker image found locally").Times(1)
		app.EXPECT().createContainer().Return(pruMock, nil)
		p.EXPECT().SetPercentage(70)
		p.EXPECT().SetState("Created Docker container")
		task.StartOnCreation = true
		app.EXPECT().StartAsync().Return(startAsyncTaskMock).Times(1)
		startAsyncTaskMock.EXPECT().GetID().Return(tskID).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		startAsyncTaskMock.EXPECT().Wait().Return(errors.New("failed to start app")).Times(1)
		app.EXPECT().SetStatus(statusFailed).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err == nil {
			t.Error("Run() should return an error when the app fails to start")
		}

		// happy case, start on creation is false, installer metadata is nil, docker image is available locally
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().GetInstaller(task.InstallerID).Return(inst, nil).Times(1)
		inst.EXPECT().GetMetadata(task.InstallerVersion).Return(core.InstallerMetadata{}, nil).Times(1)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(app, nil).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		p.EXPECT().SetPercentage(10).Times(1)
		p.EXPECT().SetState("Created application").Times(1)
		inst.EXPECT().IsPlatformImageAvailable(task.InstallerVersion).Return(true, nil).Times(1)
		p.EXPECT().SetPercentage(50).Times(1)
		p.EXPECT().SetState("Docker image found locally").Times(1)
		app.EXPECT().createContainer().Return(pruMock, nil)
		p.EXPECT().SetPercentage(70)
		p.EXPECT().SetState("Created Docker container")
		task.StartOnCreation = false
		app.EXPECT().SetStatus(statusRunning).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err != nil {
			t.Errorf("Run() should NOT return an error: %s", err.Error())
		}

		// happy case, installer metadata is available, start on creation is true, docker image is available locally
		task.InstallerMetadata = &core.InstallerMetadata{}
		amMock.EXPECT().getAppStore().Return(store).Times(1)
		store.EXPECT().CreateTemporaryInstaller(task.InstallerID, map[string]core.InstallerMetadata{task.InstallerVersion: *task.InstallerMetadata}).Return(inst)
		amMock.EXPECT().createAppForTask(task.InstallerID, task.InstallerVersion, task.AppName, task.InstallerParams, core.InstallerMetadata{}, tskID).Return(app, nil).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		p.EXPECT().SetPercentage(10).Times(1)
		p.EXPECT().SetState("Created application").Times(1)
		// docker image download
		inst.EXPECT().IsPlatformImageAvailable(task.InstallerVersion).Return(true, nil).Times(1)
		p.EXPECT().SetPercentage(50).Times(1)
		p.EXPECT().SetState("Docker image found locally").Times(1)
		// create container
		app.EXPECT().createContainer().Return(pruMock, nil)
		p.EXPECT().SetPercentage(70)
		p.EXPECT().SetState("Created Docker container")
		// start on boot
		task.StartOnCreation = true
		app.EXPECT().StartAsync().Return(startAsyncTaskMock).Times(1)
		startAsyncTaskMock.EXPECT().GetID().Return(tskID).Times(1)
		app.EXPECT().AddTask(tskID).Times(1)
		startAsyncTaskMock.EXPECT().Wait().Return(nil).Times(1)
		// set status running
		app.EXPECT().SetStatus(statusRunning).Times(1)
		err = task.Run(baseTaskMock, tskID, p)
		if err != nil {
			t.Errorf("Run() should NOT return an error: %s", err.Error())
		}

	})

	t.Run("StartAppTask", func(t *testing.T) {

		p := mock.NewMockProgress(ctrl)
		baseTaskMock := mock.NewMockTask(ctrl)
		app := NewMockapp(ctrl)
		task := StartAppTask{
			app: app,
		}

		//
		// Run
		//

		p.EXPECT().SetPercentage(50).Times(1)
		app.EXPECT().AddTask("1").Times(1)
		app.EXPECT().Start().Times(1)
		err := task.Run(baseTaskMock, "1", p)
		if err != nil {
			t.Errorf("Run() should NOT return an error: %s", err.Error())
		}
	})

	t.Run("StopAppTask", func(t *testing.T) {

		p := mock.NewMockProgress(ctrl)
		baseTaskMock := mock.NewMockTask(ctrl)
		app := NewMockapp(ctrl)
		task := StopAppTask{
			app: app,
		}

		//
		// Run
		//

		p.EXPECT().SetPercentage(50).Times(1)
		app.EXPECT().AddTask("1").Times(1)
		app.EXPECT().Stop().Times(1)
		err := task.Run(baseTaskMock, "1", p)
		if err != nil {
			t.Errorf("Run() should NOT return an error: %s", err.Error())
		}
	})

	t.Run("RemoveAppTask", func(t *testing.T) {

		p := mock.NewMockProgress(ctrl)
		baseTaskMock := mock.NewMockTask(ctrl)
		task := RemoveAppTask{
			am:    nil,
			appID: "1",
		}

		//
		// Run
		//

		// application manager is nil
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("RemoveAppTask should panic when the am field is not set")
				}
			}()
			task.Run(baseTaskMock, "1", p)
		}()

		task.am = amMock
		p.EXPECT().SetState("Deleting application").Times(1)
		p.EXPECT().SetPercentage(50).Times(1)
		amMock.EXPECT().Remove("1")
		err := task.Run(baseTaskMock, "1", p)
		if err != nil {
			t.Errorf("Run() should NOT return an error: %s", err.Error())
		}
	})

}