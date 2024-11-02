'use client'

import {useState, useCallback} from 'react'
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card"
import {AlertCircle, ExternalLink} from 'lucide-react'
import {useToast} from "@/hooks/use-toast"
import {ToastAction} from "@/components/ui/toast"

export default function LinkShortener() {
  const [longUrl, setLongUrl] = useState('')
  const [customShort, setCustomShort] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const {toast} = useToast()

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    if (!longUrl) {
      setError('Please enter a URL')
      setIsLoading(false)
      return
    }

    try {
      new URL(longUrl)
    } catch {
      setError('Please enter a valid URL')
      setIsLoading(false)
      return
    }

    const formData = new FormData()
    formData.append('original', longUrl)
    if (customShort) {
      formData.append('short', customShort)
    }

    try {
      const res = await fetch('http://localhost:1323/shorten', {
        method: 'POST',
        body: formData,
      })

      if (!res.ok) {
        throw new Error('Failed to shorten URL')
      }

      const data = await res.json()
      const shortUrl = `/${data.short}`

      toast({
        title: "URL Shortened Successfully",
        description: "Your link has been shortened. Click the button to open it.",
        action: (
          <ToastAction altText="Open shortened link" asChild>
            <Button
              variant="outline"
              size="sm"
              onClick={() => window.open(shortUrl, '_blank')}
            >
              Open Link
              <ExternalLink className="ml-2 h-4 w-4"/>
            </Button>
          </ToastAction>
        ),
      })
    } catch (err) {
      setError('An error occurred while shortening the URL')
    } finally {
      setIsLoading(false)
    }
  }, [longUrl, customShort, toast])

  return (<>
      <Card className="w-full max-w-md mx-auto">
        <CardHeader>
          <CardTitle>Link Shortener</CardTitle>
          <CardDescription>Enter a long URL to get a shortened version</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="longUrl">Long URL</Label>
              <Input
                id="longUrl"
                type="url"
                placeholder="https://example.com/very/long/url"
                value={longUrl}
                onChange={(e) => setLongUrl(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="customShort">Custom short code (optional)</Label>
              <Input
                id="customShort"
                type="text"
                placeholder="custom-code"
                value={customShort}
                onChange={(e) => setCustomShort(e.target.value)}
              />
            </div>
            {error && (
              <div className="text-red-500 flex items-center gap-2">
                <AlertCircle className="h-4 w-4"/>
                <span>{error}</span>
              </div>
            )}
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? 'Shortening...' : 'Shorten URL'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  )
}