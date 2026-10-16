param(
    [switch]$CheckOnly
)

$ErrorActionPreference = 'Stop'
$projectRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
$statePath = Join-Path $projectRoot 'CONTEXT_STATE.json'
$introPath = Join-Path $PSScriptRoot 'context_intro.md'
$outputPath = Join-Path $projectRoot 'PROJECT_CONTEXT.md'
$stateText = [System.IO.File]::ReadAllText($statePath, [System.Text.Encoding]::UTF8)
$state = $stateText | ConvertFrom-Json
$snapshotRef = $state.repository.snapshot_ref

if ($state.schema_version -ne 2) { throw 'Unsupported context schema version.' }
if ($snapshotRef -notmatch '^[a-f0-9]{40}$') { throw 'repository.snapshot_ref must be a fixed 40-character source commit SHA, not a branch or tag.' }
& git -C $projectRoot cat-file -e ($snapshotRef + '^{commit}')
if ($LASTEXITCODE -ne 0) { throw 'The fixed source commit is unavailable in this checkout.' }
$externalTargets = [System.Collections.Generic.HashSet[string]]::new()
if ($state.context_revision -lt 1) { throw 'context_revision must be positive.' }
foreach ($dateValue in @($state.context_created_on, $state.context_updated_on)) {
    [void][datetime]::ParseExact($dateValue, 'yyyy-MM-dd', [System.Globalization.CultureInfo]::InvariantCulture)
}
if ($state.user_timezone -ne 'Europe/Moscow') { throw 'Check the user timezone before generating the context.' }

function Get-NormalizedText([string]$Path) {
    return [System.IO.File]::ReadAllText($Path, [System.Text.Encoding]::UTF8).Replace("`r`n", "`n").Replace("`r", "`n").TrimEnd()
}

function Add-PortableDocument([string]$Content, [string]$SourceRelativePath, [object[]]$Documents) {
    $insideFence = $false
    $documentLines = [System.Collections.Generic.List[string]]::new()
    foreach ($line in ($Content -split "`n")) {
        if ($line -match '^\s*```') {
            $insideFence = -not $insideFence
            $documentLines.Add($line)
            continue
        }
        $portableLine = $line
        if (-not $insideFence) {
            if ($portableLine -match '^(#{1,4})\s') { $portableLine = '##' + $portableLine }
            $portableLine = [regex]::Replace($portableLine, '(?<prefix>\]\(|(?:href|src)=")(?<target>[^)"\s]+)(?<suffix>\)|")', [System.Text.RegularExpressions.MatchEvaluator]{
                param($linkMatch)
                $target = $linkMatch.Groups['target'].Value
                if ($target -match '^[a-zA-Z][a-zA-Z0-9+.-]*:' -or $target.StartsWith('#') -or $target.StartsWith('//')) { return $linkMatch.Value }
                $targetParts = $target -split '#', 2
                $sourceDirectory = [System.IO.Path]::GetDirectoryName((Join-Path $projectRoot $SourceRelativePath))
                $resolvedTarget = [System.IO.Path]::GetFullPath((Join-Path $sourceDirectory $targetParts[0]))
                foreach ($document in $Documents) {
                    $documentPath = [System.IO.Path]::GetFullPath((Join-Path $projectRoot $document.path))
                    if ($resolvedTarget.Equals($documentPath, [System.StringComparison]::OrdinalIgnoreCase) -and $targetParts.Count -eq 1) {
                        return $linkMatch.Groups['prefix'].Value + '#' + $document.anchor + $linkMatch.Groups['suffix'].Value
                    }
                }
                $rootPrefix = $projectRoot.TrimEnd('\', '/') + [System.IO.Path]::DirectorySeparatorChar
                if (-not $resolvedTarget.StartsWith($rootPrefix, [System.StringComparison]::OrdinalIgnoreCase)) { throw ('Portable link leaves project: ' + $target) }
                $repositoryRelative = $resolvedTarget.Substring($rootPrefix.Length).Replace('\', '/')
                if ($externalTargets.Add($repositoryRelative)) {
                    & git -C $projectRoot cat-file -e ($snapshotRef + ':' + $repositoryRelative)
                    if ($LASTEXITCODE -ne 0) { throw ('Snapshot resource is absent from the source commit: ' + $repositoryRelative) }
                }
                $fragment = if ($targetParts.Count -eq 2) { '#' + $targetParts[1] } else { '' }
                return $linkMatch.Groups['prefix'].Value + $state.repository.url + '/blob/' + $snapshotRef + '/' + $repositoryRelative + $fragment + $linkMatch.Groups['suffix'].Value
            })
        }
        $documentLines.Add($portableLine)
    }
    if ($insideFence) { throw 'Unbalanced code fences in a source document.' }
    return [string]::Join("`n", $documentLines)
}

$documents = @($state.source_documents)
if ($documents.Count -eq 0) { throw 'The portable context needs source documents.' }
$anchors = @($documents | ForEach-Object { $_.anchor })
if (($anchors | Select-Object -Unique).Count -ne $anchors.Count) { throw 'Duplicate appendix anchors.' }

$sourceRecords = [System.Collections.Generic.List[object]]::new()
$sourceRecords.Add([pscustomobject]@{ Path = 'CONTEXT_STATE.json'; FullPath = $statePath; Hash = (Get-FileHash -LiteralPath $statePath -Algorithm SHA256).Hash.ToLowerInvariant() })
$sourceRecords.Add([pscustomobject]@{ Path = 'tools/context_intro.md'; FullPath = $introPath; Hash = (Get-FileHash -LiteralPath $introPath -Algorithm SHA256).Hash.ToLowerInvariant() })
$sourceRecords.Add([pscustomobject]@{ Path = 'tools/update_context.ps1'; FullPath = $PSCommandPath; Hash = (Get-FileHash -LiteralPath $PSCommandPath -Algorithm SHA256).Hash.ToLowerInvariant() })
$portableDocuments = [System.Collections.Generic.List[object]]::new()

foreach ($document in $documents) {
    $sourcePath = [System.IO.Path]::GetFullPath((Join-Path $projectRoot $document.path))
    $rootPrefix = $projectRoot.TrimEnd('\', '/') + [System.IO.Path]::DirectorySeparatorChar
    if (-not $sourcePath.StartsWith($rootPrefix, [System.StringComparison]::OrdinalIgnoreCase)) {
        throw ('Source path leaves the project: ' + $document.path)
    }
    if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) { throw ('Missing source: ' + $document.path) }
    $sourceRecords.Add([pscustomobject]@{ Path = $document.path; FullPath = $sourcePath; Hash = (Get-FileHash -LiteralPath $sourcePath -Algorithm SHA256).Hash.ToLowerInvariant() })
    $sourceContent = Get-NormalizedText $sourcePath
    if ($document.PSObject.Properties['language']) {
        $sourceContent = '```' + $document.language + "`n" + $sourceContent.TrimEnd() + "`n" + '```'
    }
    $portableDocuments.Add([pscustomobject]@{ Document = $document; Content = (Add-PortableDocument $sourceContent $document.path $documents) })
}

$parts = [System.Collections.Generic.List[string]]::new()
$parts.Add('# PROJECT_CONTEXT — Faultline Lab')
$parts.Add('Самостоятельный переносимый контекст для ИИ-агентов и чатов ChatGPT. Полные источники включены ниже; локальный компьютер для чтения не требуется.')
$parts.Add(('Ревизия: **{0}**. Дата снимка: **{1}**, {2}. Исходный коммит внешних материалов: `{3}`.' -f $state.context_revision, $state.context_updated_on, $state.user_timezone, $snapshotRef))
$parts.Add('## Выполнение задач')
$parts.Add(('Актуальные статусы, блокеры и ревью ведутся только в GitHub: [Issues]({0}/issues), [Pull requests]({0}/pulls), [Milestones]({0}/milestones). Этот файл хранит спецификации, принятые решения и датированные доказательства; он не копирует текущую доску задач.' -f $state.repository.url))
$parts.Add((Get-NormalizedText $introPath))
$parts.Add('## Решения, доказательства и история')
$parts.Add('JSON ниже включён для автономного чтения решений и истории. Даты и исходные commit обязательны при использовании сохранённых наблюдений; они не означают актуальный статус задачи. Отсутствие измерения не является нулевым значением или успешной проверкой.')
$parts.Add('```json' + "`n" + $stateText.Replace("`r`n", "`n").Replace("`r", "`n").Trim() + "`n" + '```')
$parts.Add('## Реестр включённых источников')
$parts.Add('SHA-256 относится к исходному файлу на момент сборки. Генерируемый PROJECT_CONTEXT.md исключён из реестра, чтобы не создавать циклическую контрольную сумму. Проверка CheckOnly сравнивает весь снимок с результатом повторной сборки из актуальных источников.')
$registryRows = [System.Collections.Generic.List[string]]::new()
$registryRows.Add('| Источник | SHA-256 |')
$registryRows.Add('|---|---|')
foreach ($source in $sourceRecords) { $registryRows.Add(('| `{0}` | `{1}` |' -f $source.Path, $source.Hash)) }
$parts.Add([string]::Join("`n", $registryRows))
$parts.Add('## Полные самостоятельные приложения')
$indexRows = [System.Collections.Generic.List[string]]::new()
foreach ($document in $documents) { $indexRows.Add(('- [{0}](#{1}) — `{2}`' -f $document.title, $document.anchor, $document.path)) }
$parts.Add([string]::Join("`n", $indexRows))

foreach ($portableDocument in $portableDocuments) {
    $document = $portableDocument.Document
    $parts.Add(('<a id="{0}"></a>' -f $document.anchor))
    $parts.Add(('## Приложение: {0}' -f $document.title))
    $parts.Add(('Полный исходный документ: `{0}`. Текст включён непосредственно ниже; чтение исходного пути не требуется.' -f $document.path))
    $parts.Add($portableDocument.Content)
}

$generatedText = [string]::Join("`n`n", $parts) + "`n"
if ($generatedText.Contains([char]0xFFFD)) { throw 'Unexpected Unicode replacement character in context.' }

if ($CheckOnly) {
    if (-not (Test-Path -LiteralPath $outputPath -PathType Leaf)) { throw 'PROJECT_CONTEXT.md does not exist. Generate it first.' }
    $existingText = [System.IO.File]::ReadAllText($outputPath, [System.Text.Encoding]::UTF8).Replace("`r`n", "`n").Replace("`r", "`n")
    if ($existingText -cne $generatedText) { throw 'PROJECT_CONTEXT.md is stale or edited independently. Update the inputs/revision and regenerate it.' }
    Write-Output ('Context current: revision {0}; {1} embedded documents; {2} source fingerprints.' -f $state.context_revision, $documents.Count, $sourceRecords.Count)
    return
}

[System.IO.File]::WriteAllText($outputPath, $generatedText, [System.Text.UTF8Encoding]::new($false))
Write-Output ('Generated {0}; revision {1}; {2} characters; {3} embedded documents.' -f $outputPath, $state.context_revision, $generatedText.Length, $documents.Count)
