package types

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/quorum/common"
	"github.com/ethereum/quorum/crypto"
	testifyassert "github.com/stretchr/testify/assert"
)

var (
	NETWORKADMIN = "NWADMIN"
	ORGADMIN     = "OADMIN"
	NODE1        = "enode://ac6b1096ca56b9f6d004b779ae3728bf83f8e22453404cc3cef16a3d9b96608bc67c4b30db88e0a5a6c6390213f7acbe1153ff6d23ce57380104288ae19373ef@127.0.0.1:21000?discport=0&raftport=50401"
	NODE2        = "enode://0ba6b9f606a43a95edc6247cdb1c1e105145817be7bcafd6b2c0ba15d58145f0dc1a194f70ba73cd6f4cdd6864edc7687f311254c7555cc32e4d45aeb1b80416@127.0.0.1:21001?discport=0&raftport=50402"
)

var Acct1 = common.BytesToAddress([]byte("permission"))
var Acct2 = common.BytesToAddress([]byte("perm-test"))

func TestSetSyncStatus(t *testing.T) {
	assert := testifyassert.New(t)

	SetSyncStatus()

	// check if the value is set properly by calling Get
	syncStatus := GetSyncStatus()
	assert.True(syncStatus == true, fmt.Sprintf("Expected syncstatus %v . Got %v ", true, syncStatus))
}

func TestSetDefaults(t *testing.T) {
	assert := testifyassert.New(t)

	SetDefaults(NETWORKADMIN, ORGADMIN)

	// get the default values and confirm the same
	networkAdminRole, orgAdminRole, defaultAccess := GetDefaults()

	assert.True(networkAdminRole == NETWORKADMIN, fmt.Sprintf("Expected network admin role %v, got %v", NETWORKADMIN, networkAdminRole))
	assert.True(orgAdminRole == ORGADMIN, fmt.Sprintf("Expected network admin role %v, got %v", ORGADMIN, orgAdminRole))
	assert.True(defaultAccess == FullAccess, fmt.Sprintf("Expected network admin role %v, got %v", FullAccess, defaultAccess))

	SetDefaultAccess()
	networkAdminRole, orgAdminRole, defaultAccess = GetDefaults()
	assert.True(defaultAccess == ReadOnly, fmt.Sprintf("Expected network admin role %v, got %v", ReadOnly, defaultAccess))
}

func TestOrgCache_UpsertOrg(t *testing.T) {
	assert := testifyassert.New(t)

	//add a org and get the org details
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	orgInfo := OrgInfoMap.GetOrg(NETWORKADMIN)

	assert.False(orgInfo == nil, fmt.Sprintf("Expected org details, got nil"))
	assert.True(orgInfo.OrgId == NETWORKADMIN, fmt.Sprintf("Expected org id %v, got %v", NETWORKADMIN, orgInfo.OrgId))

	// update org status to suspended
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgSuspended)
	orgInfo = OrgInfoMap.GetOrg(NETWORKADMIN)

	assert.True(orgInfo.Status == OrgSuspended, fmt.Sprintf("Expected org status %v, got %v", OrgSuspended, orgInfo.Status))

	//add another org and check get org list
	OrgInfoMap.UpsertOrg(ORGADMIN, "", ORGADMIN, big.NewInt(1), OrgApproved)
	orgList := OrgInfoMap.GetOrgList()
	assert.True(len(orgList) == 2, fmt.Sprintf("Expected 2 entries, got %v", len(orgList)))

	//add sub org and check get orglist
	OrgInfoMap.UpsertOrg("SUB1", ORGADMIN, ORGADMIN, big.NewInt(2), OrgApproved)
	orgList = OrgInfoMap.GetOrgList()
	assert.True(len(orgList) == 3, fmt.Sprintf("Expected 3 entries, got %v", len(orgList)))

	//suspend the sub org and check get orglist
	OrgInfoMap.UpsertOrg("SUB1", ORGADMIN, ORGADMIN, big.NewInt(2), OrgSuspended)
	orgList = OrgInfoMap.GetOrgList()
	assert.True(len(orgList) == 3, fmt.Sprintf("Expected 3 entries, got %v", len(orgList)))
}

func TestNodeCache_UpsertNode(t *testing.T) {
	assert := testifyassert.New(t)

	// add a node into the cache and validate
	NodeInfoMap.UpsertNode(NETWORKADMIN, NODE1, NodeApproved)
	nodeInfo := NodeInfoMap.GetNodeByUrl(NODE1)
	assert.False(nodeInfo == nil, fmt.Sprintf("Expected node details, got nil"))
	assert.True(nodeInfo.OrgId == NETWORKADMIN, fmt.Sprintf("Expected org id for node %v, got %v", NETWORKADMIN, nodeInfo.OrgId))
	assert.True(nodeInfo.Url == NODE1, fmt.Sprintf("Expected node id %v, got %v", NODE1, nodeInfo.Url))

	// add another node and validate the list function
	NodeInfoMap.UpsertNode(ORGADMIN, NODE2, NodeApproved)
	nodeList := NodeInfoMap.GetNodeList()
	assert.True(len(nodeList) == 2, fmt.Sprintf("Expected 2 entries, got %v", len(nodeList)))

	// check node details update by updating node status
	NodeInfoMap.UpsertNode(ORGADMIN, NODE2, NodeDeactivated)
	nodeInfo = NodeInfoMap.GetNodeByUrl(NODE2)
	assert.True(nodeInfo.Status == NodeDeactivated, fmt.Sprintf("Expected node status %v, got %v", NodeDeactivated, nodeInfo.Status))
}

func TestRoleCache_UpsertRole(t *testing.T) {
	assert := testifyassert.New(t)

	// add a role into the cache and validate
	RoleInfoMap.UpsertRole(NETWORKADMIN, NETWORKADMIN, true, true, FullAccess, true)
	roleInfo := RoleInfoMap.GetRole(NETWORKADMIN, NETWORKADMIN)
	assert.False(roleInfo == nil, fmt.Sprintf("Expected role details, got nil"))
	assert.True(roleInfo.OrgId == NETWORKADMIN, fmt.Sprintf("Expected org id for node %v, got %v", NETWORKADMIN, roleInfo.OrgId))
	assert.True(roleInfo.RoleId == NETWORKADMIN, fmt.Sprintf("Expected node id %v, got %v", NETWORKADMIN, roleInfo.RoleId))

	// add another role and validate the list function
	RoleInfoMap.UpsertRole(ORGADMIN, ORGADMIN, true, true, FullAccess, true)
	roleList := RoleInfoMap.GetRoleList()
	assert.True(len(roleList) == 2, fmt.Sprintf("Expected 2 entries, got %v", len(roleList)))

	// update role status and validate
	RoleInfoMap.UpsertRole(ORGADMIN, ORGADMIN, true, true, FullAccess, false)
	roleInfo = RoleInfoMap.GetRole(ORGADMIN, ORGADMIN)
	assert.True(roleInfo.Active == false, fmt.Sprintf("Expected role active status to be %v, got %v", true, roleInfo.Active))
}

func TestAcctCache_UpsertAccount(t *testing.T) {
	assert := testifyassert.New(t)

	// add an account into the cache and validate
	AcctInfoMap.UpsertAccount(NETWORKADMIN, NETWORKADMIN, Acct1, true, AcctActive)
	acctInfo := AcctInfoMap.GetAccount(Acct1)
	assert.False(acctInfo == nil, fmt.Sprintf("Expected account details, got nil"))
	assert.True(acctInfo.OrgId == NETWORKADMIN, fmt.Sprintf("Expected org id for the account to be %v, got %v", NETWORKADMIN, acctInfo.OrgId))
	assert.True(acctInfo.AcctId == Acct1, fmt.Sprintf("Expected account id %x, got %x", Acct1, acctInfo.AcctId))

	// add a second account and validate the list function
	AcctInfoMap.UpsertAccount(ORGADMIN, ORGADMIN, Acct2, true, AcctActive)
	acctList := AcctInfoMap.GetAcctList()
	assert.True(len(acctList) == 2, fmt.Sprintf("Expected 2 entries, got %v", len(acctList)))

	// update account status and validate
	AcctInfoMap.UpsertAccount(ORGADMIN, ORGADMIN, Acct2, true, AcctBlacklisted)
	acctInfo = AcctInfoMap.GetAccount(Acct2)
	assert.True(acctInfo.Status == AcctBlacklisted, fmt.Sprintf("Expected account status to be %v, got %v", AcctBlacklisted, acctInfo.Status))

	// validate the list for org and role functions
	acctList = AcctInfoMap.GetAcctListOrg(NETWORKADMIN)
	assert.True(len(acctList) == 1, fmt.Sprintf("Expected number of accounts for the org to be 1, got %v", len(acctList)))
	acctList = AcctInfoMap.GetAcctListRole(NETWORKADMIN, NETWORKADMIN)
	assert.True(len(acctList) == 1, fmt.Sprintf("Expected number of accounts for the role to be 1, got %v", len(acctList)))
}

func TestGetAcctAccess(t *testing.T) {
	assert := testifyassert.New(t)

	// default access when the cache is not populated, should return default access
	SetDefaults(NETWORKADMIN, ORGADMIN)
	SetDefaultAccess()
	access := GetAcctAccess(Acct1)
	assert.True(access == ReadOnly, fmt.Sprintf("Expected account access to be %v, got %v", ReadOnly, access))

	// Create an org with two roles and two accounts linked to different roles. Validate account access
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	RoleInfoMap.UpsertRole(NETWORKADMIN, NETWORKADMIN, true, true, FullAccess, true)
	RoleInfoMap.UpsertRole(NETWORKADMIN, "ROLE1", true, true, FullAccess, true)
	AcctInfoMap.UpsertAccount(NETWORKADMIN, NETWORKADMIN, Acct1, true, AcctActive)
	AcctInfoMap.UpsertAccount(NETWORKADMIN, "ROLE1", Acct2, true, AcctActive)

	access = GetAcctAccess(Acct1)
	assert.True(access == FullAccess, fmt.Sprintf("Expected account access to be %v, got %v", FullAccess, access))

	// mark the org as pending suspension. The account access should not change
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgPendingSuspension)
	access = GetAcctAccess(Acct1)
	assert.True(access == FullAccess, fmt.Sprintf("Expected account access to be %v, got %v", FullAccess, access))

	// suspend the org and the account access should be readonly now
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgSuspended)
	access = GetAcctAccess(Acct1)
	assert.True(access == ReadOnly, fmt.Sprintf("Expected account access to be %v, got %v", ReadOnly, access))

	// mark the role as inactive and account access should now nbe read only
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	RoleInfoMap.UpsertRole(NETWORKADMIN, "ROLE1", true, true, FullAccess, false)
	access = GetAcctAccess(Acct2)
	assert.True(access == ReadOnly, fmt.Sprintf("Expected account access to be %v, got %v", ReadOnly, access))
}

func TestValidateNodeForTxn(t *testing.T) {
	assert := testifyassert.New(t)
	// pass the enode as null and the response should be true
	txnAllowed := ValidateNodeForTxn("", Acct1)
	assert.True(txnAllowed == true, "Expected access %v, got %v", true, txnAllowed)

	SetDefaultAccess()

	// if a proper enode id is not passed, return should be false
	txnAllowed = ValidateNodeForTxn("ABCDE", Acct1)
	assert.True(txnAllowed == false, "Expected access %v, got %v", true, txnAllowed)

	// if cache is not populated but the enode and account details are proper,
	// should return true
	txnAllowed = ValidateNodeForTxn(NODE1, Acct1)
	assert.True(txnAllowed == true, "Expected access %v, got %v", true, txnAllowed)

	// populate an org, account and node. validate access
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	NodeInfoMap.UpsertNode(NETWORKADMIN, NODE1, NodeApproved)
	AcctInfoMap.UpsertAccount(NETWORKADMIN, NETWORKADMIN, Acct1, true, AcctActive)
	txnAllowed = ValidateNodeForTxn(NODE1, Acct1)
	assert.True(txnAllowed == true, "Expected access %v, got %v", true, txnAllowed)

	// test access from a node not linked to the org. should return false
	OrgInfoMap.UpsertOrg(ORGADMIN, "", ORGADMIN, big.NewInt(1), OrgApproved)
	NodeInfoMap.UpsertNode(ORGADMIN, NODE2, NodeApproved)
	AcctInfoMap.UpsertAccount(ORGADMIN, ORGADMIN, Acct2, true, AcctActive)
	txnAllowed = ValidateNodeForTxn(NODE1, Acct2)
	assert.True(txnAllowed == false, "Expected access %v, got %v", true, txnAllowed)
}

// This is to make sure enode.ParseV4() honors single hexNodeId value eventhough it does follow enode URI scheme
func TestValidateNodeForTxn_whenUsingOnlyHexNodeId(t *testing.T) {
	OrgInfoMap.UpsertOrg(NETWORKADMIN, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	NodeInfoMap.UpsertNode(NETWORKADMIN, NODE1, NodeApproved)
	AcctInfoMap.UpsertAccount(NETWORKADMIN, NETWORKADMIN, Acct1, true, AcctActive)
	arbitraryPrivateKey, _ := crypto.GenerateKey()
	hexNodeId := fmt.Sprintf("%x", crypto.FromECDSAPub(&arbitraryPrivateKey.PublicKey)[1:])

	SetDefaultAccess()

	txnAllowed := ValidateNodeForTxn(hexNodeId, Acct1)

	testifyassert.False(t, txnAllowed)
}

// test the cache limit
func TestLRUCacheLimit(t *testing.T) {
	for i := 0; i < defaultOrgMapLimit ; i++ {
		orgName := "ORG" + strconv.Itoa(i)
		OrgInfoMap.UpsertOrg(orgName, "", NETWORKADMIN, big.NewInt(1), OrgApproved)
	}

	o := OrgInfoMap.GetOrg("ORG1")
	testifyassert.True(t, o != nil)
}
