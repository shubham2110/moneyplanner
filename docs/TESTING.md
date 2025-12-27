# MoneyPlanner API - Testing & Batch Scripts

Quick-reference testing scripts and batch files for Windows, Linux, and macOS.

---

## Windows PowerShell Scripts

### test-api-basic.ps1
**Initialize database and verify setup**

```powershell
# MoneyPlanner API - Basic Test
Write-Host "MoneyPlanner API - Test Suite" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

$apiUrl = "http://localhost:8080/api/init"
$headers = @{"Content-Type"="application/json"}

# Test 1: Initialize Database
Write-Host "Test 1: Initialize Database" -ForegroundColor Yellow
$body = @{force_migrate = $true} | ConvertTo-Json
try {
    $response = Invoke-RestMethod -Uri $apiUrl -Method Post -Headers $headers -Body $body
    if ($response.success) {
        Write-Host "✓ PASS - Database initialized" -ForegroundColor Green
        Write-Host "  Admin: $($response.data.admin_user.username)"
        Write-Host "  Wallet: $($response.data.default_wallet.name)"
    } else {
        Write-Host "✗ FAIL - $($response.message)" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ FAIL - Connection error: $_" -ForegroundColor Red
}

# Test 2: Check Status
Write-Host ""
Write-Host "Test 2: Check Database Status" -ForegroundColor Yellow
$body = @{force_migrate = $false} | ConvertTo-Json
try {
    $response = Invoke-RestMethod -Uri $apiUrl -Method Post -Headers $headers -Body $body
    if ($response.success) {
        Write-Host "✓ PASS - Status check successful" -ForegroundColor Green
        Write-Host "  New DB: $($response.data.database_is_new)"
        Write-Host "  Migration Required: $($response.data.migration_required)"
    } else {
        Write-Host "✗ FAIL - $($response.message)" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ FAIL - Connection error: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Tests completed!" -ForegroundColor Cyan
```

**Usage:**
```powershell
cd c:\Users\Shivani\Documents\Shubham\moneyplanner
.\test-api-basic.ps1
```

---

### test-api-full.ps1
**Comprehensive testing with custom data**

```powershell
# MoneyPlanner API - Full Test Suite
param(
    [string]$ApiUrl = "http://localhost:8080/api/init",
    [switch]$VerboseOutput
)

function Test-API {
    param([string]$TestName, [hashtable]$RequestBody)
    
    Write-Host ""
    Write-Host "Test: $TestName" -ForegroundColor Yellow
    Write-Host "-" * 50
    
    try {
        $body = $RequestBody | ConvertTo-Json
        $response = Invoke-RestMethod -Uri $ApiUrl `
          -Method Post `
          -Headers @{"Content-Type"="application/json"} `
          -Body $body
        
        if ($response.success) {
            Write-Host "✓ PASS" -ForegroundColor Green
            if ($VerboseOutput) {
                $response | ConvertTo-Json -Depth 10
            }
            return $response
        } else {
            Write-Host "✗ FAIL" -ForegroundColor Red
            Write-Host "Message: $($response.message)"
            Write-Host "Error: $($response.error)"
            if ($response.missing_items) {
                Write-Host "Missing Items:"
                $response.missing_items | ForEach-Object { Write-Host "  - $_" }
            }
            return $response
        }
    } catch {
        Write-Host "✗ ERROR: $_" -ForegroundColor Red
        return $null
    }
}

Write-Host "MoneyPlanner API - Full Test Suite" -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan

# Test 1: Basic Init
Test-API -TestName "Initialize with Defaults" -RequestBody @{
    force_migrate = $true
}

# Test 2: Custom Values
Test-API -TestName "Initialize with Custom Data" -RequestBody @{
    force_migrate = $true
    default_wallet_name = "Test Wallet"
    default_wallet_group = "TestGroup"
    admin_username = "testadmin"
    admin_password = "TestPass123"
    admin_email = "test@example.com"
    admin_name = "Test Admin"
}

# Test 3: Status Check
Test-API -TestName "Status Check (No Force)" -RequestBody @{
    force_migrate = $false
}

Write-Host ""
Write-Host "All tests completed!" -ForegroundColor Cyan
```

**Usage:**
```powershell
.\test-api-full.ps1 -VerboseOutput
```

---

### test-api-performance.ps1
**Performance and stress testing**

```powershell
# MoneyPlanner API - Performance Test
param(
    [int]$Iterations = 10,
    [int]$ConcurrentRequests = 3
)

Write-Host "MoneyPlanner API - Performance Test" -ForegroundColor Cyan
Write-Host "Iterations: $Iterations" -ForegroundColor Yellow
Write-Host "Concurrent: $ConcurrentRequests" -ForegroundColor Yellow
Write-Host ""

$apiUrl = "http://localhost:8080/api/init"
$body = @{force_migrate = $false} | ConvertTo-Json
$headers = @{"Content-Type"="application/json"}

$results = @()
$startTime = Get-Date

for ($i = 1; $i -le $Iterations; $i++) {
    $requestStart = Get-Date
    try {
        $response = Invoke-RestMethod -Uri $apiUrl `
          -Method Post `
          -Headers $headers `
          -Body $body
        
        $requestTime = ((Get-Date) - $requestStart).TotalMilliseconds
        $results += @{
            Iteration = $i
            Time = $requestTime
            Success = $response.success
        }
        
        Write-Host "[$i/$Iterations] Success in ${requestTime}ms" -ForegroundColor Green
    } catch {
        $results += @{
            Iteration = $i
            Time = -1
            Success = $false
        }
        Write-Host "[$i/$Iterations] FAILED" -ForegroundColor Red
    }
}

$totalTime = ((Get-Date) - $startTime).TotalSeconds
$successCount = ($results | Where-Object {$_.Success}).Count
$avgTime = ($results | Where-Object {$_.Time -gt 0} | Measure-Object -Property Time -Average).Average

Write-Host ""
Write-Host "Results:" -ForegroundColor Cyan
Write-Host "========" -ForegroundColor Cyan
Write-Host "Total Time: ${totalTime}s"
Write-Host "Successful: $successCount/$Iterations"
Write-Host "Average Response Time: ${avgTime}ms"
Write-Host "Min Time: $($results | Where-Object {$_.Time -gt 0} | Measure-Object -Property Time -Minimum | Select-Object -ExpandProperty Minimum)ms"
Write-Host "Max Time: $($results | Where-Object {$_.Time -gt 0} | Measure-Object -Property Time -Maximum | Select-Object -ExpandProperty Maximum)ms"
```

**Usage:**
```powershell
.\test-api-performance.ps1 -Iterations 20 -ConcurrentRequests 5
```

---

## Bash Scripts (Linux/macOS)

### test-api.sh
**Complete test suite for Unix-like systems**

```bash
#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/init"
PASS_COUNT=0
FAIL_COUNT=0

echo -e "${CYAN}MoneyPlanner API - Test Suite${NC}"
echo -e "${CYAN}=============================${NC}\n"

# Function to run test
run_test() {
    local test_name="$1"
    local payload="$2"
    
    echo -e "${YELLOW}Test: $test_name${NC}"
    
    response=$(curl -s -X POST "$API_URL" \
      -H "Content-Type: application/json" \
      -d "$payload")
    
    success=$(echo "$response" | jq '.success')
    
    if [ "$success" = "true" ]; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASS_COUNT++))
    else
        echo -e "${RED}✗ FAIL${NC}"
        echo "Message: $(echo "$response" | jq -r '.message')"
        ((FAIL_COUNT++))
    fi
    echo ""
}

# Test 1: Initialize
run_test "Initialize Database" \
  '{"force_migrate": true, "default_wallet_name": "My Wallet"}'

# Test 2: Status Check
run_test "Check Status" \
  '{"force_migrate": false}'

# Test 3: With Custom Values
run_test "Initialize with Custom Data" \
  '{"force_migrate": true, "admin_username": "testuser", "admin_email": "test@example.com"}'

# Summary
echo -e "${CYAN}Results:${NC}"
echo -e "Passed: ${GREEN}$PASS_COUNT${NC}"
echo -e "Failed: ${RED}$FAIL_COUNT${NC}"
echo -e "Total: $((PASS_COUNT + FAIL_COUNT))"
```

**Usage:**
```bash
chmod +x test-api.sh
./test-api.sh
```

---

### test-performance.sh
**Performance testing script for Unix**

```bash
#!/bin/bash

API_URL="http://localhost:8080/api/init"
ITERATIONS=${1:-10}
TOTAL_TIME=0
SUCCESS_COUNT=0

echo "MoneyPlanner API - Performance Test"
echo "Iterations: $ITERATIONS"
echo "====================================="
echo ""

for i in $(seq 1 $ITERATIONS); do
    start=$(date +%s%N)
    
    response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL" \
      -H "Content-Type: application/json" \
      -d '{"force_migrate": false}')
    
    end=$(date +%s%N)
    elapsed=$(( ($end - $start) / 1000000 ))
    
    http_code=$(echo "$response" | tail -n1)
    
    if [ "$http_code" = "200" ]; then
        ((SUCCESS_COUNT++))
        echo "[$i/$ITERATIONS] OK - ${elapsed}ms"
    else
        echo "[$i/$ITERATIONS] FAILED - HTTP $http_code"
    fi
    
    ((TOTAL_TIME += elapsed))
done

AVG_TIME=$(( TOTAL_TIME / ITERATIONS ))

echo ""
echo "Results:"
echo "--------"
echo "Successful: $SUCCESS_COUNT/$ITERATIONS"
echo "Average Time: ${AVG_TIME}ms"
echo "Total Time: ${TOTAL_TIME}ms"
```

**Usage:**
```bash
chmod +x test-performance.sh
./test-performance.sh 20
```

---

## Batch Scripts (Windows CMD)

### test-api.bat
**Simple test for Windows Command Prompt**

```batch
@echo off
REM MoneyPlanner API - Simple Test

setlocal enabledelayedexpansion

set API_URL=http://localhost:8080/api/init
set PASS_COUNT=0
set FAIL_COUNT=0

echo.
echo MoneyPlanner API - Test Suite
echo =============================
echo.

REM Test 1: Initialize
echo Test 1: Initialize Database
curl -X POST %API_URL% ^
  -H "Content-Type: application/json" ^
  -d "{\"force_migrate\": true}" ^
  -s | findstr "\"success\": true" > nul

if !errorlevel! equ 0 (
    echo [PASS] Database initialized
    set /a PASS_COUNT+=1
) else (
    echo [FAIL] Database initialization failed
    set /a FAIL_COUNT+=1
)

echo.

REM Test 2: Status Check
echo Test 2: Check Status
curl -X POST %API_URL% ^
  -H "Content-Type: application/json" ^
  -d "{\"force_migrate\": false}" ^
  -s | findstr "\"success\": true" > nul

if !errorlevel! equ 0 (
    echo [PASS] Status check successful
    set /a PASS_COUNT+=1
) else (
    echo [FAIL] Status check failed
    set /a FAIL_COUNT+=1
)

echo.
echo Results:
echo --------
echo Passed: %PASS_COUNT%
echo Failed: %FAIL_COUNT%
echo.
pause
```

**Usage:**
```batch
test-api.bat
```

---

## Docker Testing

### Dockerfile for testing
**Run tests in container**

```dockerfile
FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o app main.go

EXPOSE 8080

CMD ["./app"]
```

**Build and test:**
```bash
docker build -t moneyplanner .
docker run -p 8080:8080 moneyplanner
```

---

## Integration Test Suite (JavaScript)

### test.js
**Node.js test runner**

```javascript
const http = require('http');

const tests = [
    {
        name: 'Initialize Database',
        payload: { force_migrate: true }
    },
    {
        name: 'Check Status',
        payload: { force_migrate: false }
    },
    {
        name: 'Custom Initialization',
        payload: {
            force_migrate: true,
            admin_username: 'testuser',
            admin_email: 'test@example.com'
        }
    }
];

let passed = 0;
let failed = 0;

function runTest(test) {
    return new Promise((resolve) => {
        const data = JSON.stringify(test.payload);
        const options = {
            hostname: 'localhost',
            port: 8080,
            path: '/api/init',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': data.length
            }
        };

        const req = http.request(options, (res) => {
            let body = '';
            res.on('data', chunk => body += chunk);
            res.on('end', () => {
                try {
                    const response = JSON.parse(body);
                    if (response.success) {
                        console.log(`✓ PASS: ${test.name}`);
                        passed++;
                    } else {
                        console.log(`✗ FAIL: ${test.name}`);
                        console.log(`  Error: ${response.error}`);
                        failed++;
                    }
                } catch (e) {
                    console.log(`✗ FAIL: ${test.name} - Parse error`);
                    failed++;
                }
                resolve();
            });
        });

        req.on('error', (e) => {
            console.log(`✗ ERROR: ${test.name} - ${e.message}`);
            failed++;
            resolve();
        });

        req.write(data);
        req.end();
    });
}

async function runAllTests() {
    console.log('MoneyPlanner API - Test Suite\n');
    
    for (const test of tests) {
        await runTest(test);
    }
    
    console.log(`\nResults: ${passed} passed, ${failed} failed`);
    process.exit(failed > 0 ? 1 : 0);
}

runAllTests();
```

**Usage:**
```bash
node test.js
```

---

## Scheduled Testing (Windows Task Scheduler)

### Create Scheduled Task

```powershell
# Run tests every hour
$action = New-ScheduledTaskAction -Execute "PowerShell.exe" `
  -Argument "-File C:\path\to\test-api-basic.ps1"

$trigger = New-ScheduledTaskTrigger -RepetitionInterval "01:00:00" -RepetitionDuration "7.00:00:00"

Register-ScheduledTask -TaskName "MoneyPlanner-API-Tests" `
  -Action $action -Trigger $trigger -RunLevel Highest
```

---

## Continuous Monitoring

### monitor.ps1
**Watch API health status**

```powershell
$ApiUrl = "http://localhost:8080/api/init"
$IntervalSeconds = 5

while ($true) {
    $timestamp = Get-Date -Format "HH:mm:ss"
    
    try {
        $response = Invoke-RestMethod -Uri $ApiUrl `
          -Method Post `
          -Headers @{"Content-Type"="application/json"} `
          -Body (@{force_migrate = $false} | ConvertTo-Json)
        
        if ($response.success) {
            Write-Host "[$timestamp] ✓ API is UP" -ForegroundColor Green
        } else {
            Write-Host "[$timestamp] ✗ API returned error: $($response.message)" -ForegroundColor Red
        }
    } catch {
        Write-Host "[$timestamp] ✗ Connection failed" -ForegroundColor Red
    }
    
    Start-Sleep -Seconds $IntervalSeconds
}
```

**Run:**
```powershell
.\monitor.ps1
```

---

Last Updated: December 27, 2025
