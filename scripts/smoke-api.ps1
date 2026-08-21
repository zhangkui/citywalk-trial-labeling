$ErrorActionPreference = 'Stop'

$base = 'http://localhost:18081'
$suffix = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$adminPassword = 'Admin123!'
$adminPasswordTemp = "Admin123!$suffix"
$tempUsername = "cwtemp_$suffix"
$tempPassword1 = "Temp123!$suffix"
$tempPassword2 = "Temp456!$suffix"
$tempRoleName = "test_role_$suffix"

$counts = [ordered]@{ pass = 0; fail = 0 }

function Invoke-Api {
  param(
    [Parameter(Mandatory = $true)][string]$Method,
    [Parameter(Mandatory = $true)][string]$Path,
    [hashtable]$Body = $null,
    [string]$Token = $null,
    [string]$ForwardedFor = $null,
    [switch]$AllowHttpError
  )

  $headers = @{}
  if ($Token) { $headers.Authorization = "Bearer $Token" }
  if ($ForwardedFor) { $headers.'X-Forwarded-For' = $ForwardedFor }

  $params = @{
    Method = $Method
    Uri = "$base$Path"
    Headers = $headers
    UseBasicParsing = $true
  }
  if ($null -ne $Body) {
    $params.ContentType = 'application/json'
    $params.Body = ($Body | ConvertTo-Json -Depth 20 -Compress)
  }

  try {
    $resp = Invoke-WebRequest @params
    $json = if ($resp.Content) { $resp.Content | ConvertFrom-Json } else { $null }
    return [pscustomobject]@{ Status = [int]$resp.StatusCode; Json = $json }
  } catch {
    if (-not $AllowHttpError) { throw }
    $text = $null
    if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
      $text = $_.ErrorDetails.Message
    }
    if (-not $text) {
      $resp = $_.Exception.Response
      if ($null -eq $resp) { throw }
      $stream = $resp.GetResponseStream()
      $reader = New-Object System.IO.StreamReader($stream)
      $text = $reader.ReadToEnd()
    }
    $json = if ($text) { $text | ConvertFrom-Json } else { $null }
    $status = if ($_.Exception.Response) { [int]$_.Exception.Response.StatusCode } else { 0 }
    return [pscustomobject]@{ Status = $status; Json = $json }
  }
}

function Assert-Ok {
  param(
    [Parameter(Mandatory = $true)]$Result,
    [Parameter(Mandatory = $true)][string]$Label
  )

  if ($null -eq $Result.Json -or $Result.Json.code -ne 0) {
    throw "$Label failed: $($Result.Json | ConvertTo-Json -Depth 12)"
  }
  $script:counts.pass++
  Write-Host "PASS $Label"
  return $Result.Json.data
}

function Assert-Code {
  param(
    [Parameter(Mandatory = $true)]$Result,
    [Parameter(Mandatory = $true)][int]$Code,
    [Parameter(Mandatory = $true)][string]$Label
  )

  if ($null -eq $Result.Json -or [int]$Result.Json.code -ne $Code) {
    throw "$Label expected code $Code but got $($Result.Json.code)"
  }
  $script:counts.pass++
  Write-Host "PASS $Label"
}

function Assert-Equal {
  param(
    [Parameter(Mandatory = $true)]$Actual,
    [Parameter(Mandatory = $true)]$Expected,
    [Parameter(Mandatory = $true)][string]$Label
  )

  if ($Actual -ne $Expected) {
    throw "$Label expected $Expected but got $Actual"
  }
  $script:counts.pass++
  Write-Host "PASS $Label"
}

function Find-ByName {
  param(
    [Parameter(Mandatory = $true)]$Items,
    [Parameter(Mandatory = $true)][string]$Name
  )

  return @($Items | Where-Object { $_.name -eq $Name } | Select-Object -First 1)
}

function To-Iso {
  param([int]$DaysOffset, [int]$HoursOffset)
  return (Get-Date).ToUniversalTime().AddDays($DaysOffset).AddHours($HoursOffset).ToString('o')
}

$adminLogin = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = 'admin'; password = $adminPassword } -ForwardedFor '203.0.113.1') 'admin login'
$adminAccess = $adminLogin.accessToken
$adminRefresh = $adminLogin.refreshToken

Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/health') 'health'

$adminMe = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/me' -Token $adminAccess) 'admin me'
$adminProfileBackup = @{
  nickname = $adminMe.nickname
  avatar = $adminMe.avatar
  bio = $adminMe.bio
  city = $adminMe.city
}

Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/badges' -Token $adminAccess) 'admin badges'

$refresh = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/refresh' -Body @{ refreshToken = $adminRefresh }) 'auth refresh'
$adminAccess = $refresh.accessToken
$adminRefresh = $refresh.refreshToken

Assert-Ok (Invoke-Api -Method 'PUT' -Path '/api/v1/auth/profile' -Token $adminAccess -Body @{ nickname = 'Admin Smoke'; avatar = ''; bio = 'smoke'; city = 'Shanghai' }) 'profile update'
$profileAfter = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/me' -Token $adminAccess) 'profile verify'
Assert-Equal $profileAfter.nickname 'Admin Smoke' 'profile nickname'

Assert-Ok (Invoke-Api -Method 'PUT' -Path '/api/v1/auth/password' -Token $adminAccess -Body @{ oldPassword = $adminPassword; newPassword = $adminPasswordTemp }) 'password change temp'
$adminLoginTemp = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = 'admin'; password = $adminPasswordTemp } -ForwardedFor '203.0.113.2') 'admin relogin temp'
$adminAccess = $adminLoginTemp.accessToken
$adminRefresh = $adminLoginTemp.refreshToken

$roles = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/admin/roles' -Token $adminAccess) 'admin roles list'
$permissions = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/admin/permissions' -Token $adminAccess) 'admin permissions list'
Assert-Equal (@($permissions).Count) 29 'permission count'

$explorerRole = @($roles | Where-Object { $_.id -eq 5 } | Select-Object -First 1)
if (-not $explorerRole) { throw 'missing explorer role' }
$explorerRoleName = [string]$explorerRole.name

$tempRole = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/admin/roles' -Token $adminAccess -Body @{ name = $tempRoleName; description = 'smoke role' }) 'admin create role'
$tempRoleId = $tempRole.id
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/admin/roles/$tempRoleId" -Token $adminAccess -Body @{ name = $tempRoleName; description = 'smoke role updated' }) 'admin update role'

$tempUser = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/admin/users' -Token $adminAccess -Body @{ username = $tempUsername; password = $tempPassword1; nickname = 'Temp User'; city = 'Guangzhou'; roleIds = @() }) 'admin create user'
$tempUserId = $tempUser.id

$tempLogin1 = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = $tempUsername; password = $tempPassword1 } -ForwardedFor '203.0.113.3') 'temp login 1'
$tempAccess = $tempLogin1.accessToken
$tempRefresh = $tempLogin1.refreshToken

$tempMe1 = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/me' -Token $tempAccess) 'temp me 1'
$script:counts.pass++
Write-Host "PASS temp roles initial count=$(@($tempMe1.roles).Count)"

Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/admin/users/$tempUserId/status" -Token $adminAccess -Body @{ status = 0 }) 'disable temp user'
$disabledLogin = Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = $tempUsername; password = $tempPassword1 } -ForwardedFor '203.0.113.4' -AllowHttpError
Assert-Code $disabledLogin 1003 'disabled login blocked'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/admin/users/$tempUserId/status" -Token $adminAccess -Body @{ status = 1 }) 'enable temp user'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/admin/users/$tempUserId/password" -Token $adminAccess -Body @{ password = $tempPassword2 }) 'reset temp password'
$tempLogin2 = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = $tempUsername; password = $tempPassword2 } -ForwardedFor '203.0.113.5') 'temp login 2'
$tempAccess = $tempLogin2.accessToken
$tempRefresh = $tempLogin2.refreshToken

Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/admin/users/$tempUserId/roles" -Token $adminAccess -Body @{ roleIds = @($explorerRole.id) }) 'assign explorer role'
$tempLogin3 = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/login' -Body @{ username = $tempUsername; password = $tempPassword2 } -ForwardedFor '203.0.113.6') 'temp login 3'
$tempAccess = $tempLogin3.accessToken
$tempRefresh = $tempLogin3.refreshToken
$tempMe2 = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/me' -Token $tempAccess) 'temp me 2'
if (-not (@($tempMe2.roles) -contains $explorerRoleName)) { throw 'role assignment not reflected in token' }
$script:counts.pass++
Write-Host 'PASS temp role reflected'

$rbacDenied = Invoke-Api -Method 'GET' -Path '/api/v1/admin/users' -Token $tempAccess -AllowHttpError
Assert-Code $rbacDenied 1003 'temp denied admin users'

$routeBody = @{ 
  title = "Smoke Route $suffix"
  description = 'route smoke'
  themeId = 1
  city = 'Shanghai'
  startLat = 31.2304
  startLng = 121.4737
  endLat = 31.2354
  endLng = 121.4800
  totalDistance = 1800
  duration = 35
  difficulty = 2
  coverImage = 'https://example.com/route.jpg'
  waypoints = @(
    @{ name = 'Start'; lat = 31.2304; lng = 121.4737; stayDuration = 5; order = 1 },
    @{ name = 'Cafe'; lat = 31.2325; lng = 121.4760; stayDuration = 20; order = 2 },
    @{ name = 'End'; lat = 31.2354; lng = 121.4800; stayDuration = 10; order = 3 }
  )
}
$route = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/routes' -Token $tempAccess -Body $routeBody) 'route create'
$routeId = $route.id

$routeDetail0 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/routes/$routeId") 'route detail 0'
$routeFavoriteBefore = [int]$routeDetail0.favoriteCount
$routeRatingBefore = [int]$routeDetail0.ratingCount

Assert-Ok (Invoke-Api -Method 'POST' -Path "/api/v1/routes/$routeId/favorite" -Token $tempAccess -Body @{ }) 'route favorite'
$routeDetail1 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/routes/$routeId") 'route detail 1'
Assert-Equal ([int]$routeDetail1.favoriteCount) ($routeFavoriteBefore + 1) 'route favorite count +1'

Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/routes/$routeId/favorite" -Token $tempAccess) 'route unfavorite'
$routeDetail2 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/routes/$routeId") 'route detail 2'
Assert-Equal ([int]$routeDetail2.favoriteCount) $routeFavoriteBefore 'route favorite count reset'

Assert-Ok (Invoke-Api -Method 'POST' -Path "/api/v1/routes/$routeId/rate" -Token $tempAccess -Body @{ score = 4.5 }) 'route rate'
$routeDetail3 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/routes/$routeId") 'route detail 3'
if ([int]$routeDetail3.ratingCount -le $routeRatingBefore) { throw 'route rating count not updated' }
$script:counts.pass++
Write-Host 'PASS route rating count'

$routeUpdate = @{
  title = "Smoke Route $suffix updated"
  description = 'route smoke updated'
  themeId = 1
  city = 'Shanghai'
  startLat = 31.2304
  startLng = 121.4737
  endLat = 31.2354
  endLng = 121.4800
  totalDistance = 1900
  duration = 40
  difficulty = 2
  coverImage = 'https://example.com/route-2.jpg'
  status = 2
  waypoints = @(
    @{ name = 'Start'; lat = 31.2304; lng = 121.4737; stayDuration = 5; order = 1 },
    @{ name = 'Cafe'; lat = 31.2325; lng = 121.4760; stayDuration = 25; order = 2 },
    @{ name = 'End'; lat = 31.2354; lng = 121.4800; stayDuration = 10; order = 3 }
  )
}
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/routes/$routeId" -Token $tempAccess -Body $routeUpdate) 'route update'
$routeReviewDenied = Invoke-Api -Method 'PUT' -Path "/api/v1/routes/$routeId/status" -Token $tempAccess -Body @{ status = 2 } -AllowHttpError
Assert-Code $routeReviewDenied 1003 'route review denied'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/routes/$routeId/status" -Token $adminAccess -Body @{ status = 2 }) 'route review approve'

$comment = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/comments' -Token $tempAccess -Body @{ targetType = 'route'; targetId = $routeId; content = 'route comment'; parentId = 0 }) 'route comment create'
$commentId = $comment.id
Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/comments?targetType=route&targetId=$routeId") 'route comment list'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/comments/$commentId" -Token $tempAccess -Body @{ content = 'route comment updated' }) 'route comment update'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/comments/$commentId" -Token $tempAccess) 'route comment delete'

$tempBadges = Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/auth/badges' -Token $tempAccess) 'temp badges'
if (-not (@($tempBadges).Count -ge 1)) { throw 'badge unlock not visible after route create' }
$script:counts.pass++
Write-Host 'PASS badge unlock'

$landmark = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/landmarks' -Token $adminAccess -Body @{ name = "Smoke Landmark $suffix"; address = 'Shanghai'; lat = 31.2310; lng = 121.4780; categoryId = 1; description = 'landmark smoke'; coverImage = 'https://example.com/landmark.jpg' }) 'landmark create'
$landmarkId = $landmark.id
Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/landmarks/$landmarkId") 'landmark detail'
Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/landmarks') 'landmark list'
Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/landmarks/nearby?lat=31.2310&lng=121.4780&radius=2000&limit=10') 'landmark nearby'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/landmarks/$landmarkId" -Token $adminAccess -Body @{ name = "Smoke Landmark $suffix updated"; address = 'Shanghai'; lat = 31.2310; lng = 121.4780; categoryId = 1; description = 'landmark smoke updated'; coverImage = 'https://example.com/landmark-2.jpg' }) 'landmark update'
Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/landmarks/import' -Token $adminAccess -Body @{ }) 'landmark import'

$story = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/stories' -Token $adminAccess -Body @{ title = "Smoke Story $suffix"; content = 'story smoke'; coverImage = 'https://example.com/story.jpg'; landmarkIds = @($landmarkId); routeId = $routeId; media = @(@{ type = 1; url = 'https://example.com/story-1.jpg'; order = 1 }) }) 'story create'
$storyId = $story.id
$storyDetail0 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/stories/$storyId") 'story detail 0'
$storyLikeBefore = [int]$storyDetail0.likeCount
Assert-Ok (Invoke-Api -Method 'POST' -Path "/api/v1/stories/$storyId/like" -Token $tempAccess -Body @{ }) 'story like'
$storyDetail1 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/stories/$storyId") 'story detail 1'
Assert-Equal ([int]$storyDetail1.likeCount) ($storyLikeBefore + 1) 'story like count +1'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/stories/$storyId/like" -Token $tempAccess) 'story unlike'
$storyDetail2 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/stories/$storyId") 'story detail 2'
Assert-Equal ([int]$storyDetail2.likeCount) $storyLikeBefore 'story like count reset'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/stories/$storyId" -Token $adminAccess -Body @{ title = "Smoke Story $suffix updated"; content = 'story smoke updated'; coverImage = 'https://example.com/story-2.jpg'; landmarkIds = @($landmarkId); routeId = $routeId; status = 2; media = @(@{ type = 1; url = 'https://example.com/story-2.jpg'; order = 1 }) }) 'story update'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/stories/$storyId/status" -Token $adminAccess -Body @{ status = 2 }) 'story review approve'

$landmarkComment = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/comments' -Token $tempAccess -Body @{ targetType = 'landmark'; targetId = $landmarkId; content = 'landmark comment'; parentId = 0 }) 'landmark comment create'
$landmarkCommentId = $landmarkComment.id
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/comments/$landmarkCommentId" -Token $tempAccess) 'landmark comment delete'

$eventMeetup = To-Iso -DaysOffset 1 -HoursOffset 1
$eventStart = To-Iso -DaysOffset 1 -HoursOffset 2
$eventEnd = To-Iso -DaysOffset 1 -HoursOffset 3
$event = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/events' -Token $adminAccess -Body @{ title = "Smoke Event $suffix"; routeId = $routeId; meetupAddress = 'Shanghai'; meetupLat = 31.2310; meetupLng = 121.4780; meetupTime = $eventMeetup; startTime = $eventStart; endTime = $eventEnd; maxParticipants = 20; fee = 0; description = 'event smoke' }) 'event create'
$eventId = $event.id
$eventDetail0 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/events/$eventId") 'event detail 0'
$participantsBefore = [int]$eventDetail0.currentParticipants
Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/events') 'event list'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/events/$eventId" -Token $adminAccess -Body @{ title = "Smoke Event $suffix updated"; routeId = $routeId; meetupAddress = 'Shanghai'; meetupLat = 31.2310; meetupLng = 121.4780; meetupTime = $eventMeetup; startTime = $eventStart; endTime = $eventEnd; maxParticipants = 25; fee = 100; description = 'event smoke updated'; status = 1 }) 'event update'
$join = Assert-Ok (Invoke-Api -Method 'POST' -Path "/api/v1/events/$eventId/join" -Token $tempAccess -Body @{ name = 'Temp User'; remark = 'smoke' }) 'event join'
$eventDetail1 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/events/$eventId") 'event detail 1'
Assert-Equal ([int]$eventDetail1.currentParticipants) ($participantsBefore + 1) 'event participant +1'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/events/$eventId/checkin" -Token $tempAccess -Body @{ }) 'event checkin'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/events/$eventId/join" -Token $tempAccess) 'event cancel join'
$eventDetail2 = Assert-Ok (Invoke-Api -Method 'GET' -Path "/api/v1/events/$eventId") 'event detail 2'
Assert-Equal ([int]$eventDetail2.currentParticipants) $participantsBefore 'event participant reset'

Assert-Ok (Invoke-Api -Method 'GET' -Path '/api/v1/recommend?lat=31.2304&lng=121.4737') 'recommend'

$logout = Assert-Ok (Invoke-Api -Method 'POST' -Path '/api/v1/auth/logout' -Token $tempAccess -Body @{ refreshToken = $tempRefresh }) 'temp logout'
$tempMeAfterLogout = Invoke-Api -Method 'GET' -Path '/api/v1/auth/me' -Token $tempAccess -AllowHttpError
Assert-Code $tempMeAfterLogout 1002 'logout blacklists access token'

# Cleanup while admin token is still valid.
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/events/$eventId" -Token $adminAccess) 'event delete'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/stories/$storyId" -Token $adminAccess) 'story delete'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/landmarks/$landmarkId" -Token $adminAccess) 'landmark delete'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/routes/$routeId" -Token $adminAccess) 'route delete'
Assert-Ok (Invoke-Api -Method 'DELETE' -Path "/api/v1/admin/roles/$tempRoleId" -Token $adminAccess) 'role delete'
Assert-Ok (Invoke-Api -Method 'PUT' -Path "/api/v1/auth/profile" -Token $adminAccess -Body $adminProfileBackup) 'profile restore'
Assert-Ok (Invoke-Api -Method 'PUT' -Path '/api/v1/auth/password' -Token $adminAccess -Body @{ oldPassword = $adminPasswordTemp; newPassword = $adminPassword }) 'password restore'

$adminLogout = Invoke-Api -Method 'POST' -Path '/api/v1/auth/logout' -Token $adminAccess -Body @{ refreshToken = $adminRefresh }
Assert-Ok $adminLogout 'admin logout'

Write-Host "DONE pass=$($counts.pass) fail=$($counts.fail)"
